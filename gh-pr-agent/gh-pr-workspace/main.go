package main

import (
	"dagger/gh-pr-workspace/internal/dagger"
	"strings"
)

const GHCLIImage = "maniator/gh:v2.68.1"

type GhPrWorkspace struct {
	// The repository URL to run against
	RepositoryURL string

	// The pull request number to use
	//+internal-use-only
	PullRequestNumber int

	// The workspace container tool to use the GH CLI
	Container *dagger.Container
}

func New(repository string, pullRequestNumber int, token *dagger.Secret) *GhPrWorkspace {
	return &GhPrWorkspace{
		RepositoryURL:     strings.TrimPrefix(strings.TrimPrefix(repository, "https://"), "github.com/"),
		PullRequestNumber: pullRequestNumber,
		Container:         dag.Container().From(GHCLIImage).WithSecretVariable("GH_TOKEN", token),
	}
}
