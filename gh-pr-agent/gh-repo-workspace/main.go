package main

import (
	"dagger/gh-repo-workspace/internal/dagger"
	"strings"
)

const GHCLIImage = "maniator/gh:v2.68.1"

type GhRepoWorkspace struct{
	// The repository URL to run against
	//+internal-use-only
	URL string

	// The workspace container tool to use the GH CLI
	Container *dagger.Container

	// The current user
	//+internal-use-only
	User string
}

// Create a new repository workspace with the given repository URL and token.
func New(repository string, token *dagger.Secret) *GhRepoWorkspace {
	return &GhRepoWorkspace{
		URL:        parseRepoURL(repository),
		Container:  dag.Container().From(GHCLIImage).WithSecretVariable("GH_TOKEN", token),
		User: "",
	}
}

// WithUser sets the current user to use by the agent
// This can be useful to improve the precision of the agent's answers
// and returns more user's focused contents.
func (g *GhRepoWorkspace) WithUser(user string) *GhRepoWorkspace {
	g.User = user
	return g
}

func parseRepoURL(repository string) string {
	return  strings.TrimPrefix(strings.TrimPrefix(repository, "https://"), "github.com/")
}