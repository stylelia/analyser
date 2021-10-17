package analyser

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Cookstyle structs
type CookstyleMetadata struct {
	Version string `json:"version"`
}

// Cookstyle check payload
type CookstyleCheck struct {
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

func runCookstyle(exec CommandRunner) (CookstyleCheck, error) {
	var c CookstyleCheck

	// cmd := exec.Command("cookstyle", "-a", "--format", "json")
	output, err := exec.Output()
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(output, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func getLatestCookstyle(cookstyleApi string, client *http.Client) (string, error) {
	request, err := http.NewRequest(http.MethodGet, cookstyleApi, nil)
	if err != nil {
		return "", err
	}

	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	var getCookstyleVersion CookstyleMetadata
	err = json.NewDecoder(response.Body).Decode(&getCookstyleVersion)
	if err != nil {
		return "", err
	}

	return getCookstyleVersion.Version, nil
}

func (c *CookstyleCheck) PrintMessage(cookstyleVersion string) string {
	header := fmt.Sprintf("Hi!\n\nI ran Cookstyle %s against this repo and here are the results.\n\nSummary:\nOffence Count: %v\n\nChanges:", cookstyleVersion, c.Summary.OffenseCount)

	var logs string
	for _, part := range c.Files {
		var partial string
		if len(part.Offenses) > 0 {
			partial += fmt.Sprintf("\nIssue found and resolved with %s\n\n", part.Path)
			for _, offenses := range part.Offenses {
				partial += fmt.Sprintf("- %s\n", offenses.Message)
			}
		}
		logs += partial
	}

	return header + logs
}
