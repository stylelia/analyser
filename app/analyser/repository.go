package analyser

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Repository structs
type Repository struct {
	Org           string
	Name          string
	DefaultBranch string
	LatestCommit  string
	GithubApi     string
}

type GetRepository struct {
	DefaultBranch string `json:"default_branch"`
}

type GetLastCommit struct {
	Sha string `json:"sha"`
}

func NewRepo(org, name, defaultBranch string) Repository {
	return Repository{
		Org:           org,
		Name:          name,
		DefaultBranch: defaultBranch,
	}
}

func getDefaultBranch(endpoint string, client *http.Client) (string, error) {
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var getRepoStruct GetRepository
	err = json.NewDecoder(response.Body).Decode(&getRepoStruct)
	if err != nil {
		return "", err
	}

	return getRepoStruct.DefaultBranch, nil
}

func (r *Repository) getLastCommit(endpoint string, client *http.Client) error {
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
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

func (r *Repository) buildCommitEndpoint(githubApi string) string {
	return fmt.Sprintf("%s/repos/%s/%s/commits/%s", githubApi, r.Org, r.Name, r.DefaultBranch)
}

func (r *Repository) Clone(exec CommandRunner) error {
	err := exec.Run()
	if err != nil {
		fmt.Print(err.Error())
		return err
	}

	return nil
}
