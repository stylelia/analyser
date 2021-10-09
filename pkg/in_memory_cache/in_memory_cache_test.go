package in_memory_cache

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ctx context.Context = context.Background()

func TestUpdateSha(t *testing.T) {
	t.Run("Updates a Key which does not exist", func(t *testing.T) {
		sha1 := "b64d5bae3cee6da8c305c0f46f678914cb22e483"
		githubOrg := "stylelia"
		repoName := "updateKeyRepo"
		expected := sha1

		imc := InMemoryCache{}
		err := imc.UpdateCommitSha(ctx, githubOrg, repoName, sha1)
		assert.NoError(t, err)

		actual, err := imc.GetCommitSha(ctx, githubOrg, repoName)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

	})

	t.Run("Updates a Key which already exists", func(t *testing.T) {
		sha1 := "b64d5bae3cee6da8c305c0f46f678914cb22e483"
		sha1Updated := "b64d5bae3cee6da8c305c0f46f678914cb22e600" // End changed to 600.

		githubOrg := "stylelia"
		repoName := "updateKeyRepo"
		expected := sha1Updated

		imc := InMemoryCache{}
		err := imc.UpdateCommitSha(ctx, githubOrg, repoName, sha1)
		assert.NoError(t, err)

		err = imc.UpdateCommitSha(ctx, githubOrg, repoName, sha1Updated)
		assert.NoError(t, err)

		actual, err := imc.GetCommitSha(ctx, githubOrg, repoName)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

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
