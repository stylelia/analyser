package analyser

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cookstyleJSON CookstyleCheck = CookstyleCheck{
	Metadata: Metadata{},
	Files: []Files{
		Files{
			Path: "/tmp/path",
			Offenses: []Offenses{
				Offenses{
					Severity:    "High",
					Message:     "First message",
					Correctable: true,
				},
			},
		},
		Files{
			Path: "/tmp/another",
			Offenses: []Offenses{
				Offenses{
					Severity:    "Medium",
					Message:     "Second message",
					Correctable: true,
				},
			},
		},
	},
	Summary: Summary{
		OffenseCount: 2,
	},
}

type MockRunCookstyleCommand struct{}

func (m *MockRunCookstyleCommand) Run() error {
	return nil
}

func (m *MockRunCookstyleCommand) Output() ([]byte, error) {
	out, err := json.Marshal(cookstyleJSON)
	if err != nil {
		// we should, nor we never will panic here
		panic(err)
	}
	return out, nil
}

type MockRunCookstyleCommand_Error struct{}

func (m *MockRunCookstyleCommand_Error) Run() error {
	return errors.New("test error")
}

func (m MockRunCookstyleCommand_Error) Output() ([]byte, error) {
	return nil, errors.New("test error")
}

func TestRunCookstyle(t *testing.T) {
	t.Run("runCookstyle throws an error on a faulty command", func(t *testing.T) {
		faulty := &MockRunCookstyleCommand_Error{}

		_, err := runCookstyle(faulty)
		assert.Error(t, err)
	})

	t.Run("runCookstyle doesn't return any error on a valid command and returns a valid JSON", func(t *testing.T) {
		runner := &MockRunCookstyleCommand{}

		out, err := runCookstyle(runner)
		assert.NoError(t, err)
		assert.Equal(t, cookstyleJSON, out)
	})
}

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
		_, err = w.Write(response)
		if err != nil {
			// we should never panic here.
			// we should never reach here.
			// we shouldn't use panic.
			panic(err)
		}
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
