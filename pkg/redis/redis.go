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

func NewRedis(port uint16, server, password string) *Redis {
	Addr := fmt.Sprintf("%v:%v", server, port)
	client := redis.NewClient(&redis.Options{
		Addr:     Addr,
		Password: password, // no password set
		DB:       0,        // use default DB
	})
	return &Redis{client: client}
}

func (r *Redis) UpdateCommitSha(ctx context.Context, githubOrg, repoName, commitSha string) error {
	keyPath := r.keyPath(githubOrg, repoName)
	return r.updateKeyField(ctx, keyPath, commitShaFieldName, commitSha)
}

func (r *Redis) GetCommitSha(ctx context.Context, githubOrg, repoName string) (string, error) {
	keyPath := r.keyPath(githubOrg, repoName)
	return r.getKeyField(ctx, keyPath, commitShaFieldName)
}

func (r *Redis) UpdateToolVersion(ctx context.Context, githubOrg, repoName, toolName, toolVersion string) error {
	keyPath := r.keyPath(githubOrg, repoName)
	return r.updateKeyField(ctx, keyPath, toolName, toolVersion)
}

func (r *Redis) GetToolVersion(ctx context.Context, githubOrg, repoName, toolName string) (string, error) {
	keyPath := r.keyPath(githubOrg, repoName)
	return r.getKeyField(ctx, keyPath, toolName)
}

func (r *Redis) KeyNotFoundInCacheError() error {
	return errors.New("cache: key not found")
}

// Private methods
func (r *Redis) keyPath(githubOrg, repoName string) string {
	return fmt.Sprintf("github/%v/%v", githubOrg, repoName)
}

func (r *Redis) getKeyField(ctx context.Context, keyPath, fieldName string) (string, error) {
	value, err := r.client.HGet(ctx, keyPath, fieldName).Result()
	err = r.normaliseErrorCode(err)
	return value, err
}

func (r *Redis) updateKeyField(ctx context.Context, keyPath, fieldName, fieldValue string) error {
	err := r.client.HSet(ctx, keyPath, fieldName, fieldValue).Err()
	return r.normaliseErrorCode(err)
}

// Used to reconcile errors into error codes we know that the client can then validate against
func (r *Redis) normaliseErrorCode(err error) error {
	if err != nil && err.Error() == redis.Nil.Error() {
		err = r.KeyNotFoundInCacheError()
	}
	return err
}

func (r *Redis) deleteKey(ctx context.Context, githubOrg, repoName string) error {
	keyPath := r.keyPath(githubOrg, repoName)
	err := r.client.Del(ctx, keyPath).Err()
	return r.normaliseErrorCode(err)
}
