package analyser

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRepo(t *testing.T) {
	client := &http.Client{}

	t.Run("GetRepo fetches the Master branch on https://github.com/youshy/ENAE-The-System", func(t *testing.T) {
		org := "youshy"
		name := "ENAE-The-System"
		expected := "master"

		repo, err := getRepo(org, name, client)
		assert.NoError(t, err)
		assert.Equal(t, expected, repo.DefaultBranch)
	})

	t.Run("GetRepo fetches the Main branch on https://github.com/sous-chefs/golang", func(t *testing.T) {
		org := "sous-chefs"
		name := "golang"
		expected := "main"

		repo, err := getRepo(org, name, client)
		assert.NoError(t, err)
		assert.Equal(t, expected, repo.DefaultBranch)
	})
}

func TestGetLastCommit(t *testing.T) {
	client := &http.Client{}

	t.Run("GetLastCommit fetches the last commit on Main branch on https://github.com/sous-chefs/golang", func(t *testing.T) {
		org := "sous-chefs"
		name := "golang"

		repo, err := getRepo(org, name, client)
		assert.NoError(t, err)

		err = repo.getLastCommit(client)
		assert.NoError(t, err)
	})
}
