package redis

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	redisHost     string          = os.Getenv("REDIS_HOST")
	redisPort     uint16          = convertRedisPort(os.Getenv("REDIS_PORT"))
	redisPassword string          = os.Getenv("REDIS_PASSWORD")
	ctx           context.Context = context.Background()
)

func convertRedisPort(port string) uint16 {
	value, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		// We should never panic here...
		log.Fatalf("Error in converting Redis Port: %v\n", err)
	}
	return uint16(value)
}

func TestConvertRedisPort(t *testing.T) {
	expected := uint16(1234)
	actual := convertRedisPort("1234")
	assert.Equal(t, expected, actual)
}

func TestUpdateSha(t *testing.T) {
	t.Run("Updates a Key which does not exist", func(t *testing.T) {
		sha1 := "b64d5bae3cee6da8c305c0f46f678914cb22e483"
		githubOrg := "stylelia"
		repoName := "newKeyRepo"
		expected := sha1

		r := NewRedis(redisPort, redisHost, redisPassword)
		defer r.DeleteKey(ctx, githubOrg, repoName)
		err := r.UpdateCommitSha(ctx, githubOrg, repoName, sha1)
		assert.NoError(t, err)

		actual, err := r.GetCommitSha(ctx, githubOrg, repoName)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("Updates a Key which already exists", func(t *testing.T) {
		sha1 := "b64d5bae3cee6da8c305c0f46f678914cb22e483"
		sha1Updated := "b64d5bae3cee6da8c305c0f46f678914cb22e600" // End changed to 600.

		githubOrg := "stylelia"
		repoName := "updateKeyRepo"
		expected := sha1Updated

		r := NewRedisCache(redisPort, redisHost, redisPassword)
		defer r.DeleteKey(ctx, githubOrg, repoName)
		err := r.UpdateCommitSha(ctx, githubOrg, repoName, sha1)
		assert.NoError(t, err)

		err = r.UpdateCommitSha(ctx, githubOrg, repoName, sha1Updated)
		assert.NoError(t, err)

		actual, err := r.GetCommitSha(ctx, githubOrg, repoName)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)

	})
}

func TestGetSha(t *testing.T) {
	t.Run("Errors when trying to get a key which does not exist", func(t *testing.T) {
		githubOrg := "stylelia"
		repoName := "newKeyRepo"
		expected := ""

		r := NewRedisCache(redisPort, redisHost, redisPassword)
		actual, err := r.GetCommitSha(ctx, githubOrg, repoName)
		assert.EqualError(t, err, r.KeyNotFoundInCacheError().Error())
		assert.Equal(t, expected, actual)
	})
}
