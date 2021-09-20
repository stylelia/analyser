package analyser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

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

// Endpoints
const (
	mainApi string = "https://api.github.com"
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
