package redis_cache

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var commitShaFieldName = "commitSha"

// TODO: Keys will be a composite - read Redis docs!
// We should link to the docs here so people don't need to google it.
// Redis setup for now
type RedisCache struct {
	client *redis.Client
}

func (r *RedisCache) UpdateCommitSha(ctx context.Context, githubOrg, repoName, commitSha string) error {
	keyPath := r.keyPath(githubOrg, repoName)
	return r.client.HSet(ctx, keyPath, commitShaFieldName, commitSha).Err()
}

func (r *RedisCache) GetCommitSha(ctx context.Context, githubOrg, repoName string) (string, error) {
	keyPath := r.keyPath(githubOrg, repoName)
	value, err := r.client.HGet(ctx, keyPath, commitShaFieldName).Result()
	if err != nil && err.Error() == redis.Nil.Error() {
		err = r.KeyNotFoundInCacheError()
	}
	return value, err
}

func (r *RedisCache) DeleteKey(ctx context.Context, githubOrg, repoName string) error {
	keyPath := r.keyPath(githubOrg, repoName)
	return r.client.Del(ctx, keyPath).Err()
}

func (r *RedisCache) keyPath(githubOrg, repoName string) string {
	return fmt.Sprintf("github/%v/%v", githubOrg, repoName)
}

func (r *RedisCache) KeyNotFoundInCacheError() error {
	return errors.New("cache: key not found")
}

func NewRedisCache(port uint16, server, password string) *RedisCache {
	client := newClient(port, server, password)
	return &RedisCache{client: client}
}

func newClient(port uint16, server, password string) *redis.Client {
	Addr := fmt.Sprintf("%v:%v", server, port)
	client := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: password, // no password set
		DB:       0,        // use default DB
	})
	return client
}
