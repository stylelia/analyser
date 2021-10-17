package analyser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLatestCookstyle(t *testing.T) {
	expectedVersion := "v.10.10.10"

	cookstyleHandler := func(w http.ResponseWriter, r *http.Request) {
		output := CookstyleMetadata{
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

	server := httptest.NewServer(http.HandlerFunc(cookstyleHandler))
	defer server.Close()

	client := &http.Client{}

	ver, err := getLatestCookstyle(server.URL, client)
	assert.NoError(t, err)
	assert.Equal(t, expectedVersion, ver)
}

func TestPrintMessage(t *testing.T) {
	cookstyleVersion := "v10.2.10"

	cookstyleJSON := CookstyleCheck{
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

	validMessage := fmt.Sprintf("Hi!\n\nI ran Cookstyle %s against this repo and here are the results.\n\nSummary:\nOffence Count: %v\n\nChanges:\nIssue found and resolved with %s\n\n- %s\n\nIssue found and resolved with %s\n\n- %s\n", cookstyleVersion, cookstyleJSON.Summary.OffenseCount, cookstyleJSON.Files[0].Path, cookstyleJSON.Files[0].Offenses[0].Message, cookstyleJSON.Files[1].Path, cookstyleJSON.Files[1].Offenses[0].Message)
	t.Run("Print message returns a valid message", func(t *testing.T) {
		out := cookstyleJSON.PrintMessage(cookstyleVersion)
		assert.Equal(t, validMessage, out)
	})
}
