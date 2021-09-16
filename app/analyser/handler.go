package analyser

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

	// Check cache for cookstyle for a given repo.
	// If exists, check version - if equal and if commit sha equal to cache, leave app

	// If not exists or version is different or sha is different, clone the repo

	// run 'cookstyle -a --format json'

	// If cookstyle finds a change, create a new branch 'styleila/cookstyle_<version>'
	// If no change, update cache with cookstyle and default branch sha

	// Raise a PR for that change
	// put in pr body nice message based on json response from cookstyle

	// update cache with default branch sha & cookstyle version

	// see: https://github.com/Xorima/github-cookstyle-runner/blob/main/app/entrypoint.ps1#L139 to L157

	return nil
}

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

type GetCookbook struct {
	Version string `json:"version"`
}

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

	log.Printf("Version: %+v", getCookbookVersion)

	return getCookbookVersion.Version, nil
}
