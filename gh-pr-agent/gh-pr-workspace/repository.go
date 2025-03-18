package main

import (
	"context"
	"dagger/gh-pr-workspace/internal/dagger"
	"fmt"
)

type Repository struct {
	// The container containing the repository pull request
	//+internal-use-only
	Container *dagger.Container
}

// Repository returns a container containing the repository pull request directory.
func (g *GhPrWorkspace) Repository(ctx context.Context) (*Repository, error) {
	repo := g.Container.
		WithExec([]string{"gh", "repo", "clone", g.RepositoryURL, "/repo"}).
		WithWorkdir("/repo").
		WithExec([]string{"gh", "pr", "checkout", fmt.Sprintf("%d", g.PullRequestNumber)})

	return &Repository{
		Container: repo,
	}, nil
}

// Diff returns the diff of the pull request and the origin
func (r *Repository) Diff(ctx context.Context) (string, error) {
	return r.Container.WithExec([]string{"gh", "pr", "diff"}).Stdout(ctx)
}

// Read returns a file from the repository
func (r *Repository) Read(path string) *dagger.File {
	return r.Container.File(path)
}
