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
	redisHost     string          = getenv("REDIS_HOST", "redis")
	redisPort     uint16          = convertRedisPort(getenv("REDIS_PORT", "6379"))
	redisPassword string          = os.Getenv("REDIS_PASSWORD")
	ctx           context.Context = context.Background()
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// This is in here as it is only used in tests to clean up after,
// This is not great, but the linter was complaining a lot
func (r *Redis) deleteKey(ctx context.Context, githubOrg, repoName string) {
	keyPath := r.keyPath(githubOrg, repoName)
	err := r.client.Del(ctx, keyPath).Err()
	// We are not going to check this error in each test,
	// so let's just go bang here...
	if err != nil {
		log.Fatal(err)
	}
}

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
	t.Run("Updates a Sha Field which does not exist", func(t *testing.T) {
		sha1 := "b64d5bae3cee6da8c305c0f46f678914cb22e483"
		githubOrg := "stylelia"
		repoName := "newKeyRepo"
		expected := sha1

		r := NewRedis(redisPort, redisHost, redisPassword)
		defer r.deleteKey(ctx, githubOrg, repoName)
		err := r.UpdateCommitSha(ctx, githubOrg, repoName, sha1)
		assert.NoError(t, err)

		actual, err := r.GetCommitSha(ctx, githubOrg, repoName)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
	t.Run("Updates a Sha Field which already exists", func(t *testing.T) {
		sha1 := "b64d5bae3cee6da8c305c0f46f678914cb22e483"
		sha1Expected := "b64d5bae3cee6da8c305c0f46f678914cb22e600" // End changed to 600.

		githubOrg := "stylelia"
		repoName := "updateKeyRepo"

		r := NewRedis(redisPort, redisHost, redisPassword)
		defer r.deleteKey(ctx, githubOrg, repoName)
		err := r.UpdateCommitSha(ctx, githubOrg, repoName, sha1)
		assert.NoError(t, err)

		err = r.UpdateCommitSha(ctx, githubOrg, repoName, sha1Expected)
		assert.NoError(t, err)

		actual, err := r.GetCommitSha(ctx, githubOrg, repoName)
		assert.NoError(t, err)
		assert.Equal(t, sha1Expected, actual)

	})
}

func TestGetSha(t *testing.T) {
	t.Run("Returns no error when trying to get a Sha Field which does not exist", func(t *testing.T) {
		githubOrg := "stylelia"
		repoName := "newKeyRepo"
		expected := ""

		r := NewRedis(redisPort, redisHost, redisPassword)
		actual, err := r.GetCommitSha(ctx, githubOrg, repoName)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestUpdateToolVersion(t *testing.T) {
	t.Run("Updates a Tool with a version which does not exist", func(t *testing.T) {
		toolName := "cookstyle"
		expectedToolVersion := "1.2.3"
		githubOrg := "stylelia"
		repoName := "updateToolRepo"

		r := NewRedis(redisPort, redisHost, redisPassword)
		defer r.deleteKey(ctx, githubOrg, repoName)
		err := r.UpdateToolVersion(ctx, githubOrg, repoName, toolName, expectedToolVersion)
		assert.NoError(t, err)

		actual, err := r.GetToolVersion(ctx, githubOrg, repoName, toolName)
		assert.NoError(t, err)
		assert.Equal(t, expectedToolVersion, actual)
	})
	t.Run("Updates a Tool Field which already exists", func(t *testing.T) {
		toolName := "cookstyle"
		toolVersion := "1.0.0"
		expectedToolVersion := "1.2.3"

		githubOrg := "stylelia"
		repoName := "updateKeyRepo"

		r := NewRedis(redisPort, redisHost, redisPassword)
		defer r.deleteKey(ctx, githubOrg, repoName)
		err := r.UpdateToolVersion(ctx, githubOrg, repoName, toolName, toolVersion)
		assert.NoError(t, err)

		err = r.UpdateToolVersion(ctx, githubOrg, repoName, toolName, expectedToolVersion)
		assert.NoError(t, err)

		actual, err := r.GetToolVersion(ctx, githubOrg, repoName, toolName)
		assert.NoError(t, err)
		assert.Equal(t, expectedToolVersion, actual)

	})
}

func TestGetToolVersion(t *testing.T) {
	t.Run("Returns no error when trying to get a Tool Field which does not exist", func(t *testing.T) {
		githubOrg := "stylelia"
		repoName := "newKeyRepo"
		expected := ""
		toolName := "cookstyle"

		r := NewRedis(redisPort, redisHost, redisPassword)
		actual, err := r.GetToolVersion(ctx, githubOrg, repoName, toolName)
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}
