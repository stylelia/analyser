package analyser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
