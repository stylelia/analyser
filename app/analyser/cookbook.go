package analyser

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
)

const (
	cookbookApi string = "https://rubygems.org/api/v1/versions/cookstyle/latest.json"
)

// Cookbook structs
type GetCookbook struct {
	Version string `json:"version"`
}

// Cookbook check payload
type CookbookCheck struct {
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

func runCookbook() (CookbookCheck, error) {
	var (
		c      CookbookCheck
		cmdOut bytes.Buffer
	)

	cmd := exec.Command("cookstyle", "-a", "format", "json")
	cmd.Stdout = &cmdOut

	err := cmd.Run()
	if err != nil {
		return c, err
	}

	err = gob.NewDecoder(&cmdOut).Decode(&c)
	if err != nil {
		return c, err
	}

	return c, nil
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

	return getCookbookVersion.Version, nil
}

func (c *CookbookCheck) PrintMessage(cookstyleVersion string) string {
	var (
		logs string
	)
	header := fmt.Sprintf("Hi!\n\nI ran Cookstyle %s against this repo and here are the results.\n\nSummary:\nOffence Count: %v\n\nChanges:\n", cookstyleVersion, c.Summary.OffenseCount)

	for _, part := range c.Files {
		var partial string
		for _, offenses := range part.Offenses {
			partial += fmt.Sprintf("Issue found and resolved with %s\n\n%s\n\n", part.Path, offenses.Message)
		}
		logs += partial
	}

	return header + logs
}
