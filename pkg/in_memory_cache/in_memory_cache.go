package in_memory_cache

import (
	"context"
	"errors"
	"fmt"
)

type InMemoryCache struct{}
type cacheEntry struct {
	commitSha string
}

var (
	cache map[string]cacheEntry = make(map[string]cacheEntry)
)

func (i *InMemoryCache) UpdateCommitSha(ctx context.Context, githubOrg, repoName, commitSha string) error {
	keyPath := i.keyPath(githubOrg, repoName)
	cache[keyPath] = cacheEntry{commitSha: commitSha}
	return nil
}

func (i *InMemoryCache) GetCommitSha(ctx context.Context, githubOrg, repoName string) (string, error) {
	keyPath := i.keyPath(githubOrg, repoName)
	sha := cache[keyPath].commitSha
	if sha == "" {
		return sha, i.KeyNotFoundInCacheError()
	}
	return sha, nil
}

func (i *InMemoryCache) keyPath(githubOrg, repoName string) string {
	return fmt.Sprintf("%v/%v", githubOrg, repoName)
}

func (i *InMemoryCache) KeyNotFoundInCacheError() error {
	return errors.New("cache: key not found")
}
