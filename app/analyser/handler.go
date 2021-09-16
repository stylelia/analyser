package analyser

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleEvent() error {
	return handle()
}

func handle() error {
	// Fetch the latest default commit sha and check it against cache

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

	Client *http.Client
}

type GetRepository struct {
	DefaultBranch string `json:"default_branch"`
}

type GetLastCommit struct {
	Sha string `json:"sha"`
}

const (
	mainApi string = "https://api.github.com"
)

func getRepo(org, name string) (Repository, error) {
	// GET /repos/{owner}/{repo}
	var r Repository

	getRepoUri := fmt.Sprintf("%s/repos/%s/%s", mainApi, org, name)

	request, err := http.NewRequest(http.MethodGet, getRepoUri, nil)
	if err != nil {
		return r, err
	}

	// NOTE: This can probably go outside
	// We might have to set up oAuth2 here + timeout
	client := &http.Client{}

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

	r.Client = client

	return r, nil
}

func (r *Repository) getLastCommit() error {
	// GET /repos/:owner/:repo/commits/:branch
	getLastCommitUri := fmt.Sprintf("%s/repos/%s/%s/commits/%s", mainApi, r.Org, r.Name, r.DefaultBranch)

	request, err := http.NewRequest(http.MethodGet, getLastCommitUri, nil)
	if err != nil {
		return err
	}

	response, err := r.Client.Do(request)
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
