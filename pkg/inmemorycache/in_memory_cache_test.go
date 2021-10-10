package inmemorycache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ctx context.Context = context.Background()

func TestUpdateSha(t *testing.T) {
	t.Run("Updates a Key which does not exist", func(t *testing.T) {
		expectedSha1 := "b64d5bae3cee6da8c305c0f46f678914cb22e483"
		githubOrg := "stylelia"
		repoName := "updateKeyRepo"

		imc := NewInMemoryCache()
		err := imc.UpdateCommitSha(ctx, githubOrg, repoName, expectedSha1)
		assert.NoError(t, err)

		actual, err := imc.GetCommitSha(ctx, githubOrg, repoName)
		assert.NoError(t, err)
		assert.Equal(t, expectedSha1, actual)

	})

	t.Run("Updates a Key which already exists", func(t *testing.T) {
		sha1 := "b64d5bae3cee6da8c305c0f46f678914cb22e483"
		expectedSha1 := "b64d5bae3cee6da8c305c0f46f678914cb22e600" // End changed to 600.

		githubOrg := "stylelia"
		repoName := "updateKeyRepo"

		imc := NewInMemoryCache()
		err := imc.UpdateCommitSha(ctx, githubOrg, repoName, sha1)
		assert.NoError(t, err)

		err = imc.UpdateCommitSha(ctx, githubOrg, repoName, expectedSha1)
		assert.NoError(t, err)

		actual, err := imc.GetCommitSha(ctx, githubOrg, repoName)
		assert.NoError(t, err)
		assert.Equal(t, expectedSha1, actual)

	})
}

func TestGetSha(t *testing.T) {
	t.Run("Errors when trying to get a key which does not exist", func(t *testing.T) {
		githubOrg := "stylelia"
		repoName := "newKeyRepo"
		expected := ""

		imc := InMemoryCache{}
		actual, err := imc.GetCommitSha(ctx, githubOrg, repoName)
		assert.EqualError(t, err, imc.KeyNotFoundInCacheError().Error())
		assert.Equal(t, expected, actual)
	})
}

func TestUpdateToolVersion(t *testing.T) {
	t.Run("Updates a Tool with a version which does not exist", func(t *testing.T) {
		toolName := "cookstyle"
		expectedToolVersion := "1.2.3"
		githubOrg := "stylelia"
		repoName := "updateToolRepo"

		imc := NewInMemoryCache()
		err := imc.UpdateToolVersion(ctx, githubOrg, repoName, toolName, expectedToolVersion)
		assert.NoError(t, err)

		actual, err := imc.GetToolVersion(ctx, githubOrg, repoName, toolName)
		assert.NoError(t, err)
		assert.Equal(t, expectedToolVersion, actual)
	})
	t.Run("Updates a Tool Field which already exists", func(t *testing.T) {
		toolName := "cookstyle"
		toolVersion := "1.0.0"
		expectedToolVersion := "1.2.3"

		githubOrg := "stylelia"
		repoName := "updateKeyRepo"

		imc := NewInMemoryCache()
		err := imc.UpdateToolVersion(ctx, githubOrg, repoName, toolName, toolVersion)
		assert.NoError(t, err)

		err = imc.UpdateToolVersion(ctx, githubOrg, repoName, toolName, expectedToolVersion)
		assert.NoError(t, err)

		actual, err := imc.GetToolVersion(ctx, githubOrg, repoName, toolName)
		assert.NoError(t, err)
		assert.Equal(t, expectedToolVersion, actual)

	})
}

func TestGetToolVersion(t *testing.T) {
	t.Run("Errors when trying to get a Tool Field which does not exist", func(t *testing.T) {
		githubOrg := "stylelia"
		repoName := "newKeyRepo"
		expected := ""
		toolName := "cookstyle"

		imc := NewInMemoryCache()
		actual, err := imc.GetToolVersion(ctx, githubOrg, repoName, toolName)
		assert.EqualError(t, err, imc.KeyNotFoundInCacheError().Error())
		assert.Equal(t, expected, actual)
	})
}
