package inmemorycache

import (
	"context"
	"errors"
	"fmt"
)

type InMemoryCache struct {
	cache map[string]cacheEntry
}
type cacheEntry struct {
	commitSha string
}

func (i *InMemoryCache) UpdateCommitSha(ctx context.Context, githubOrg, repoName, commitSha string) error {
	keyPath := i.keyPath(githubOrg, repoName)
	i.cache[keyPath] = cacheEntry{commitSha: commitSha}
	return nil
}

func (i *InMemoryCache) GetCommitSha(ctx context.Context, githubOrg, repoName string) (string, error) {
	keyPath := i.keyPath(githubOrg, repoName)
	sha := i.cache[keyPath].commitSha
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

func NewInMemoryCache(port uint16, server, password string) *InMemoryCache {
	cache := make(map[string]cacheEntry)
	return &InMemoryCache{cache: cache}
}
