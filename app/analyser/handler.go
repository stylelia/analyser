package analyser

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/styleila/analyser/pkg/redis"
)

const (
	Commit    string = "Commit"
	Cookstyle string = "Cookstyle"

	githubApi    string = "https://api.github.com"
	cookstyleApi string = "https://rubygems.org/api/v1/versions/cookstyle/latest.json"
)

// Interface for KV store
type KeyValueStore interface {
	GetKey(key string) (string, error)
	UpdateKey(key, value string) error
}

// Interface for all exec.Command stuff
type CommandRunner interface {
	Run() error
	Output() ([]byte, error)
}

type Handler struct {
	Client *http.Client
}

func NewHandler(client *http.Client) Handler {
	return Handler{
		Client: client,
	}
}

func HandleEvent() error {
	client := &http.Client{}

	handler := NewHandler(client)

	return handler.handle()
}

func (h *Handler) handle() error {
	// Fetch the latest default commit sha and check it against cache
	// TODO: Setup cache
	org := os.Getenv("ORGANISATION")
	name := os.Getenv("NAME")

	githubDefaultBranchEndpoint := fmt.Sprintf("%s/repos/%s/%s", githubApi, org, name)

	branch, err := getDefaultBranch(githubDefaultBranchEndpoint, h.Client)
	if err != nil {
		return err
	}

	repo := NewRepo(org, name, branch)

	githubLastCommitEndpoint := repo.buildCommitEndpoint(githubApi)

	err = repo.getLastCommit(githubLastCommitEndpoint, h.Client)
	if err != nil {
		return err
	}

	// Setup redis
	redis := redis.Redis{}

	latestCommit, err := redis.GetKey(Commit)
	if err != nil {
		return err
	}

	// Check cache for cookstyle for a given repo.
	// If exists, check version - if equal and if commit sha equal to cache, leave app
	cookstyleVersion, err := getLatestCookstyle(cookstyleApi, h.Client)
	if err != nil {
		return err
	}

	latestCookstyle, err := redis.GetKey(Cookstyle)
	if err != nil {
		return err
	}

	if repo.LatestCommit == latestCommit && cookstyleVersion == latestCookstyle {
		// log that we're ending the lifecycle here
		return nil
	}

	// If not exists or version is different or sha is different, clone the repo
	// TODO: This requires an access to a valid SSH key on the lambda
	repoUri := fmt.Sprintf("git@github.com:%s/%s.git", repo.Org, repo.Name)
	cloneRepoRunner := exec.Command("git", "clone", repoUri)
	err = repo.clone(cloneRepoRunner)
	if err != nil {
		return err
	}

	// run 'cookstyle -a --format json'
	runner := exec.Command("cookstyle", "-a", "--format", "json")
	out, err := runCookstyle(runner)
	if err != nil {
		return err
	}

	// If cookstyle finds a change, create a new branch 'styleila/cookstyle_<version>'
	branchRunner := buildBranchCommand(createBranchName(cookstyleVersion))
	if out.Summary.OffenseCount > 0 {
		err = createBranch(branchRunner)
		if err != nil {
			return err
		}
	}

	// Raise a PR for that change
	// put in pr body nice message based on json response from cookstyle
	// TODO: Write a printer for PR

	// update cache with default branch sha & cookstyle version
	err = redis.UpdateKey(Commit, repo.LatestCommit)
	if err != nil {
		return err
	}

	err = redis.UpdateKey(Cookstyle, latestCookstyle)
	if err != nil {
		return err
	}

	// see: https://github.com/Xorima/github-cookstyle-runner/blob/main/app/entrypoint.ps1#L139 to L157

	return nil
}
