package main

import (
	"context"
	"dagger/gh-review-agent/internal/dagger"
)

type GhReviewAgent struct {
	Repo              string
	PullRequestNumber int
	Token             *dagger.Secret
}

func New(repo string, token *dagger.Secret) *GhReviewAgent {
	return &GhReviewAgent{
		Repo:  repo,
		Token: token,
	}
}

func (g *GhReviewAgent) SetPRNumber(prNumber int) *GhReviewAgent {
	g.PullRequestNumber = prNumber
	return g
}

func (g *GhReviewAgent) Ask(
	ctx context.Context,
	question string,

	//+optional
	prNumber int,
) (string, error) {
	if prNumber == 0 {
		prNumber = g.PullRequestNumber
	}

	return dag.LLM().
		SetGhPrWorkspace("gh-pr-workspace", dag.GhPrWorkspace(g.Repo, prNumber, g.Token)).
		SetGhRepoWorkspace("gh-repo-workspace", dag.GhRepoWorkspace(g.Repo, g.Token)).
		WithPromptVar("question", question).
		WithPrompt(
			`You are a helpful assistant that can answer question regarding a given pull request.

You have been given access to a workspace gr-pr-workspace containing two different tools: 
- a tool conversation to get all the messages sent in the PR including reviews and comments.
- a tool repository to get the differences between the PR and the origin or read files contents.

You have been given access to a workspace gh-repo-workspace containing on tool:
- a tool list-pull-requests to get all the pull requests in the repository with the given filters.

Answer to the given question using the tools in an efficient way.

<question>
$question
</question>
		`).
		LastReply(ctx)
}
