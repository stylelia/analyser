package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

const commitShaFieldName string = "commitSha"

// Creates a Redis client with methods for Updating and Getting the relevant keys
// That matter within the Stylelia application
// This package is opionated on purpose
type Redis struct {
	client *redis.Client
}

func (r *Redis) UpdateCommitSha(ctx context.Context, githubOrg, repoName, commitSha string) error {
	keyPath := r.keyPath(githubOrg, repoName)
	return r.client.HSet(ctx, keyPath, commitShaFieldName, commitSha).Err()
}

func (r *Redis) GetCommitSha(ctx context.Context, githubOrg, repoName string) (string, error) {
	keyPath := r.keyPath(githubOrg, repoName)
	value, err := r.client.HGet(ctx, keyPath, commitShaFieldName).Result()
	if err != nil && err.Error() == redis.Nil.Error() {
		err = r.KeyNotFoundInCacheError()
	}
	return value, err
}

func (r *Redis) DeleteKey(ctx context.Context, githubOrg, repoName string) error {
	keyPath := r.keyPath(githubOrg, repoName)
	return r.client.Del(ctx, keyPath).Err()
}

func (r *Redis) keyPath(githubOrg, repoName string) string {
	return fmt.Sprintf("github/%v/%v", githubOrg, repoName)
}

func (r *Redis) KeyNotFoundInCacheError() error {
	return errors.New("cache: key not found")
}

func NewRedis(port uint16, server, password string) *Redis {
	Addr := fmt.Sprintf("%v:%v", server, port)
	client := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: password, // no password set
		DB:       0,        // use default DB
	})
	return &Redis{client: client}
}
