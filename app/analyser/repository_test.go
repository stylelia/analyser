package analyser

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDefaultBranch(t *testing.T) {
	defaultBranch := "ObiWanKenobiHadTheHigherGround"

	gitBranchHandler := func(w http.ResponseWriter, r *http.Request) {
		output := GetRepository{
			DefaultBranch: defaultBranch,
		}
		response, err := json.Marshal(output)
		if err != nil {
			// we should never panic here.
			// we should never reach here.
			// we shouldn't use panic.
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}

	server := httptest.NewServer(http.HandlerFunc(gitBranchHandler))
	defer server.Close()

	client := &http.Client{}

	branch, err := getDefaultBranch(server.URL, client)
	assert.NoError(t, err)
	assert.Equal(t, defaultBranch, branch)
}

func TestGetLastCommit(t *testing.T) {
	latestCommit := "YoungSkywalkerWasDoomedToFail"

	gitBranchHandler := func(w http.ResponseWriter, r *http.Request) {
		output := GetLastCommit{
			Sha: latestCommit,
		}
		response, err := json.Marshal(output)
		if err != nil {
			// we should never panic here.
			// we should never reach here.
			// we shouldn't use panic.
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}

	server := httptest.NewServer(http.HandlerFunc(gitBranchHandler))
	defer server.Close()

	client := &http.Client{}

	repo := NewRepo("someOrg", "someName", "defaultBranch")

	err := repo.getLastCommit(server.URL, client)
	assert.NoError(t, err)
	assert.Equal(t, latestCommit, repo.LatestCommit)
}

func TestBuildCommitEndpoint(t *testing.T) {
	githubApi := "https://api.github.io"
	expected := "https://api.github.io/repos/org/name/commits/branch"

	repo := NewRepo("org", "name", "branch")
	endpoint := repo.buildCommitEndpoint(githubApi)

	assert.Equal(t, expected, endpoint)
}
