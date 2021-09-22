package analyser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestCookbook(t *testing.T) {
	expectedVersion := "v.10.10.10"

	cookbookHandler := func(w http.ResponseWriter, r *http.Request) {
		output := GetCookbook{
			Version: expectedVersion,
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

	server := httptest.NewServer(http.HandlerFunc(cookbookHandler))
	defer server.Close()

	client := &http.Client{}

	ver, err := getLatestCookbook(server.URL, client)
	assert.NoError(t, err)
	assert.Equal(t, expectedVersion, ver)
}

func TestPrintMessage(t *testing.T) {
	cookstyleVersion := "v10.2.10"

	cookbookJSON := CookbookCheck{
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

	validMessage := fmt.Sprintf("Hi!\n\nI ran Cookstyle %s against this repo and here are the results.\n\nSummary:\nOffence Count: %v\n\nChanges:\nIssue found and resolved with %s\n\n%s\n\nIssue found and resolved with %s\n\n%s\n\n", cookstyleVersion, cookbookJSON.Summary.OffenseCount, cookbookJSON.Files[0].Path, cookbookJSON.Files[0].Offenses[0].Message, cookbookJSON.Files[1].Path, cookbookJSON.Files[1].Offenses[0].Message)

	t.Run("Print message returns a valid message", func(t *testing.T) {
		out := cookbookJSON.PrintMessage(cookstyleVersion)
		assert.Equal(t, validMessage, out)
	})
}
