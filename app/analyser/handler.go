package analyser

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

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

	repo, err := getRepo(org, name, h.Client)
	if err != nil {
		return err
	}

	err = repo.getLastCommit(h.Client)
	if err != nil {
		return err
	}

	// TODO: Fetch latest commit from cache and test against

	// Check cache for cookstyle for a given repo.
	// If exists, check version - if equal and if commit sha equal to cache, leave app
	ver, err := getLatestCookbook(h.Client)
	if err != nil {
		return err
	}

	// TODO: Fetch latest cookstyle version from cache and test against
	log.Printf("%v\n", ver) // compiler moans

	// If not exists or version is different or sha is different, clone the repo
	// TODO: This requires an access to a valid SSH key on the lambda
	err = repo.clone()
	if err != nil {
		return err
	}

	// run 'cookstyle -a --format json'
	out, err := runCookbook()
	if err != nil {
		return err
	}

	// If cookstyle finds a change, create a new branch 'styleila/cookstyle_<version>'
	if out.Summary.OffenseCount > 0 {
		err = createBranch(ver)
		if err != nil {
			return err
		}
	}
	// If no change, update cache with cookstyle and default branch sha
	// TODO: this can be performed when cache is ready

	// Raise a PR for that change
	// put in pr body nice message based on json response from cookstyle

	// update cache with default branch sha & cookstyle version

	// see: https://github.com/Xorima/github-cookstyle-runner/blob/main/app/entrypoint.ps1#L139 to L157

	return nil
}

// Repository structs
type Repository struct {
	Org           string
	Name          string
	DefaultBranch string
	LatestCommit  string
}

type GetRepository struct {
	DefaultBranch string `json:"default_branch"`
}

type GetLastCommit struct {
	Sha string `json:"sha"`
}

// Cookbook structs
type GetCookbook struct {
	Version string `json:"version"`
}

// Cookbook check payload
type CookbookCheck struct {
	Metadata Metadata `json:"metadata"`
	Files    []Files  `json:"files"`
	Summary  Summary  `json:"summary"`
}

type Metadata struct {
	RubocopVersion string `json:"rubocop_version"`
	RubyEngine     string `json:"ruby_engine"`
	RubyVersion    string `json:"ruby_version"`
	RubyPatchlevel string `json:"ruby_patchlevel"`
	RubyPlatform   string `json:"ruby_platform"`
}

type Files struct {
	Path     string     `json:"path"`
	Offenses []Offenses `json:"offenses"`
}

type Offenses struct {
	Severity    string `json:"severity"`
	Message     string `json:"message"`
	CopName     string `json:"cop_name"`
	Corrected   bool   `json:"corrected"`
	Correctable bool   `json:"correctable"`
}

type Summary struct {
	OffenseCount       int `json:"offense_count"`
	TargetFileCount    int `json:"target_file_count"`
	InspectedFileCount int `json:"inspected_file_count"`
}

// Endpoints
const (
	mainApi     string = "https://api.github.com"
	cookbookApi string = "https://rubygems.org/api/v1/versions/cookstyle/latest.json"
)

func getRepo(org, name string, client *http.Client) (Repository, error) {
	// GET /repos/{owner}/{repo}
	var r Repository

	getRepoUri := fmt.Sprintf("%s/repos/%s/%s", mainApi, org, name)

	request, err := http.NewRequest(http.MethodGet, getRepoUri, nil)
	if err != nil {
		return r, err
	}

	response, err := client.Do(request)
	if err != nil {
		return r, err
	}
	defer response.Body.Close()

	var getRepoStruct GetRepository
	err = json.NewDecoder(response.Body).Decode(&getRepoStruct)

	r.Org = org
	r.Name = name
	r.DefaultBranch = getRepoStruct.DefaultBranch

	return r, nil
}

func (r *Repository) getLastCommit(client *http.Client) error {
	// GET /repos/:owner/:repo/commits/:branch
	getLastCommitUri := fmt.Sprintf("%s/repos/%s/%s/commits/%s", mainApi, r.Org, r.Name, r.DefaultBranch)

	request, err := http.NewRequest(http.MethodGet, getLastCommitUri, nil)
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var getLastCommitStruct GetLastCommit
	err = json.NewDecoder(response.Body).Decode(&getLastCommitStruct)
	if err != nil {
		return err
	}

	r.LatestCommit = getLastCommitStruct.Sha

	return nil
}

func (r *Repository) clone() error {
	repoUri := fmt.Sprintf("git@github.com:%s/%s.git", r.Org, r.Name)

	err := exec.Command("git", "clone", repoUri).Run()
	if err != nil {
		return err
	}

	return nil
}

func createBranch(cookbookVersion string) error {
	cmdMessage := fmt.Sprintf("stylelia/cookstyle_%s", cookbookVersion)

	err := exec.Command("git", "branch", "-b", cmdMessage).Run()
	if err != nil {
		return err
	}

	return nil
}

func runCookbook() (CookbookCheck, error) {
	var (
		c      CookbookCheck
		cmdOut bytes.Buffer
	)

	cmd := exec.Command("cookstyle", "-a", "format", "json")
	cmd.Stdout = &cmdOut

	err := cmd.Run()
	if err != nil {
		return c, err
	}

	err = gob.NewDecoder(&cmdOut).Decode(&c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func getLatestCookbook(client *http.Client) (string, error) {
	request, err := http.NewRequest(http.MethodGet, cookbookApi, nil)
	if err != nil {
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var getCookbookVersion GetCookbook
	err = json.NewDecoder(response.Body).Decode(&getCookbookVersion)
	if err != nil {
		return "", err
	}

	return getCookbookVersion.Version, nil
}
