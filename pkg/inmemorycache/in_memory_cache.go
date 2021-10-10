package inmemorycache

import (
	"context"
	"errors"
	"fmt"
)

type InMemoryCache struct {
	cache map[string]cacheEntry
}

// TODO: Make the cacheEntry support multiple tools for encase
// a repo has a mixture of languages/tools they want scanned
// currently we only support 1 tool, cookstyle, so this will
// be a future itteration
type cacheEntry struct {
	commitSha   string
	toolName    string
	toolVersion string
}

func NewInMemoryCache() *InMemoryCache {
	cache := make(map[string]cacheEntry)
	return &InMemoryCache{cache: cache}
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

func (i *InMemoryCache) KeyNotFoundInCacheError() error {
	return errors.New("cache: key not found")
}

func (i *InMemoryCache) UpdateToolVersion(ctx context.Context, githubOrg, repoName, toolName, toolVersion string) error {
	keyPath := i.keyPath(githubOrg, repoName)
	i.cache[keyPath] = cacheEntry{toolName: toolName, toolVersion: toolVersion}
	return nil
}

func (i *InMemoryCache) GetToolVersion(ctx context.Context, githubOrg, repoName, toolName string) (string, error) {
	keyPath := i.keyPath(githubOrg, repoName)
	toolVersion := i.cache[keyPath].toolVersion
	if toolVersion == "" {
		return toolVersion, i.KeyNotFoundInCacheError()
	}
	return toolVersion, nil
}

// Private methods

func (i *InMemoryCache) keyPath(githubOrg, repoName string) string {
	return fmt.Sprintf("%v/%v", githubOrg, repoName)
}
