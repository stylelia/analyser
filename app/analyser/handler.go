package analyser

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/google/go-github/v39/github"
	"github.com/styleila/analyser/pkg/redis"
	"github.com/youshy/logger"
	"go.uber.org/zap"
)

const (
	Commit       string = "Commit"
	Cookstyle    string = "Cookstyle"
	WorkingDir   string = "/tmp" // Only wriable location in lambda
	githubApi    string = "https://api.github.com"
	cookstyleApi string = "https://rubygems.org/api/v1/versions/cookstyle/latest.json"
)

// Interface for KV store
type KeyValueStore interface {
	GetCommitSha(context.Context, string, string) (string, error)
	UpdateCommitSha(context.Context, string, string, string) error
	GetToolVersion(context.Context, string, string, string) (string, error)
	UpdateToolVersion(context.Context, string, string, string, string) error
}

type Handler struct {
	Client *http.Client
	Log    *zap.SugaredLogger
}

func NewHandler(client *http.Client, log *zap.SugaredLogger) Handler {
	return Handler{
		Client: client,
		Log:    log,
	}
}

func HandleEvent() error {
	client := &http.Client{}
	// TODO: make as env flags
	logger := logger.NewLogger(logger.DEBUG, false)

	handler := NewHandler(client, logger)

	return handler.handle()
}

func (h *Handler) handle() error {
	// Fetch the latest default commit sha and check it against cache
	org := os.Getenv("ORGANISATION")
	name := os.Getenv("NAME")

	githubDefaultBranchEndpoint := fmt.Sprintf("%s/repos/%s/%s", githubApi, org, name)

	branch, err := getDefaultBranch(githubDefaultBranchEndpoint, h.Client)
	if err != nil {
		h.Log.Errorf("Unable to get default branch: %v", err)
		return err
	}

	repo := NewRepo(org, name, branch)

	githubLastCommitEndpoint := repo.buildCommitEndpoint(githubApi)

	err = repo.getLastCommit(githubLastCommitEndpoint, h.Client)
	if err != nil {
		h.Log.Errorf("Unable to get latest commit: %v", err)
		return err
	}

	// Setup redis
	portRaw := os.Getenv("REDIS_PORT")
	server := os.Getenv("REDIS_HOST")
	password := os.Getenv("REDIS_PASSWORD")

	port, err := strconv.Atoi(portRaw)
	if err != nil {
		h.Log.Errorf("Unable to convert port: %v", err)
		return err
	}

	redis := redis.NewRedis(uint16(port), server, password)
	ctx := context.Background()

	latestCommit, err := redis.GetCommitSha(ctx, org, name)
	if err != nil {
		h.Log.Errorf("Unable to get commit sha from Redis: %v", err)
		return err
	}

	// Check cache for cookstyle for a given repo.
	// If exists, check version - if equal and if commit sha equal to cache, leave app
	cookstyleVersion, err := getLatestCookstyle(cookstyleApi, h.Client)
	if err != nil {
		h.Log.Errorf("Unable to get latest cookstyle version: %v", err)
		return err
	}

	latestCookstyle, err := redis.GetToolVersion(ctx, org, name, Cookstyle)
	if err != nil {
		h.Log.Errorf("Unable to get latest cookstyle version from Redis: %v", err)
		return err
	}

	if repo.LatestCommit == latestCommit && cookstyleVersion == latestCookstyle {
		// log that we're ending the lifecycle here
		h.Log.Info("All up to date!")
		return nil
	}

	// If not exists or version is different or sha is different, clone the repo
	repoUri := fmt.Sprintf("https://%s@github.com/%s/%s.git", os.Getenv("GITHUB_TOKEN"), repo.Org, repo.Name)
	cloneRepoRunner := exec.Command("git", "clone", repoUri, WorkingDir)
	cloneRepoRunner.Dir = WorkingDir
	err = repo.Clone(cloneRepoRunner)
	if err != nil {
		h.Log.Errorf("Unable to clone repo: %v", err)
		return err
	}

	// run 'cookstyle -a --format json'
	runner := exec.Command("cookstyle", "-a", "--format", "json")
	runner.Dir = WorkingDir
	out, err := runCookstyle(runner)
	if err != nil {
		h.Log.Errorf("Unable to run cookstyle: %v", err)
		return err
	}

	// If cookstyle finds a change, create a new branch 'styleila/cookstyle_<version>'
	branchName := createBranchName(cookstyleVersion)
	branchRunner := buildBranchCommand(branchName)
	branchRunner.Dir = WorkingDir
	title := fmt.Sprintf("Stylelia: Cookstyle %s updates", cookstyleVersion)
	message := out.PrintMessage(cookstyleVersion)

	if out.Summary.OffenseCount > 0 {
		err = gitCmdRunner(branchRunner)
		if err != nil {
			h.Log.Errorf("Unable to add new branch: %v", err)
			return err
		}

		stageRunner := buildStageCommand()
		stageRunner.Dir = WorkingDir
		err = gitCmdRunner(stageRunner)
		if err != nil {
			h.Log.Errorf("Unable to stage commit: %v", err)
			return err
		}

		commitRunner := buildCommitCommand(os.Getenv("GIT_EMAIL"), os.Getenv("GIT_USERNAME"), title, message)
		commitRunner.Dir = WorkingDir
		err = gitCmdRunner(commitRunner)
		if err != nil {
			h.Log.Errorf("Unable to commit: %v", err)
			return err
		}

		pushRunner := buildPushCommand(branchName)
		pushRunner.Dir = WorkingDir
		err = gitCmdRunner(pushRunner)
		if err != nil {
			h.Log.Errorf("Unable to push commit: %v", err)
			return err
		}

		// Raise a PR for that change if one does not exist
		// put in pr body nice message based on json response from cookstyle
		client := createClientWithAuth(ctx)

		opt := &github.PullRequestListOptions{Head: branchName, State: "open"}
		existingPr, _, err := client.PullRequests.List(ctx, repo.Org, repo.Name, opt)
		if err != nil {
			h.Log.Errorf("Unable to get PRs: %v", err)
			return err
		}

		if len(existingPr) != 1 {
			pr := &github.NewPullRequest{
				Title:               &title,
				Head:                &branchName,
				Base:                &repo.DefaultBranch,
				Body:                &message,
				MaintainerCanModify: github.Bool(true),
			}

			_, _, err = client.PullRequests.Create(ctx, repo.Org, repo.Name, pr)
			if err != nil {
				h.Log.Errorf("Unable to create PR: %v", err)
				return err
			}

		} else if len(existingPr) == 1 {
			// Update body as there is some change on the PR we should reflect in the text
			existingPr[0].Body = &message
			updatedPr := existingPr[0]
			updatedPr.Body = &message
			client.PullRequests.Edit(ctx, repo.Org, repo.Name, existingPr[0].GetNumber(), updatedPr)
		}
	}

	// update cache with default branch sha & cookstyle version
	err = redis.UpdateCommitSha(ctx, org, name, repo.LatestCommit)
	if err != nil {
		h.Log.Errorf("Unable to update commit sha in Redis: %v", err)
		return err
	}

	err = redis.UpdateToolVersion(ctx, org, name, Cookstyle, cookstyleVersion)
	if err != nil {
		h.Log.Errorf("Unable to update tool version in Redis: %v", err)
		return err
	}

	return nil
}
