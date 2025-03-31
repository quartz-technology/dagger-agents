package main

import (
	"context"
	"dagger/gh-repo-workspace/internal/dagger"
	"encoding/json"
	"fmt"
)

type pullRequest struct {
	Number             int      `json:"number"`
	Title              string   `json:"title"`
	Author             string   `json:"author"`
	CreatedAt          string   `json:"created_at"`
	State              string   `json:"state"`
	Draft              bool     `json:"draft"`
	Tags               []string `json:"tags"`
	RequestedReviewers []string `json:"requested_reviewers"`
}

// Custom UnmarshalJSON function for PullRequest
// +internal-use-only
func (p *pullRequest) UnmarshalJSON(data []byte) error {
	var temp struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		User   struct {
			Login string `json:"login"`
		} `json:"user"`
		CreatedAt string `json:"created_at"`
		State     string `json:"state"`
		Draft     bool   `json:"draft"`
		Labels    []struct {
			Name string `json:"name"`
		} `json:"labels"`
		RequestedTeams []struct {
			Name string `json:"name"`
		} `json:"requested_teams"`
		RequestedReviewers []struct {
			Login string `json:"login"`
		} `json:"requested_reviewers"`
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	p.Number = temp.Number
	p.Title = temp.Title
	p.Author = temp.User.Login
	p.CreatedAt = temp.CreatedAt
	p.State = temp.State
	p.Draft = temp.Draft
	p.Tags = make([]string, len(temp.Labels))
	for i, label := range temp.Labels {
		p.Tags[i] = label.Name
	}

	p.RequestedReviewers = make([]string, len(temp.RequestedReviewers)+len(temp.RequestedTeams))
	for i, reviewer := range temp.RequestedReviewers {
		p.RequestedReviewers[i] = reviewer.Login
	}
	for i, team := range temp.RequestedTeams {
		p.RequestedReviewers[i+len(temp.RequestedReviewers)] = team.Name
	}

	return nil
}

// List Pull requests returns a file containing all the pull requests in the repository
// with the given filter applied.
func (r *GhRepoWorkspace) ListPullRequests(
	ctx context.Context,

	// The filters to apply to the pull requests, in the form of a HTML query string
	// Example: `state=open&draft=false` to list all open PRs that are not in draft mode
	//+default="state=open&draft=false"
	filters string,
) (*dagger.File, error) {
	result, err := r.Container.WithExec(
		[]string{"gh", "api", fmt.Sprintf("repos/%s/pulls?%s", r.URL, filters)},
		dagger.ContainerWithExecOpts{
			RedirectStdout: "/pull-requests.json",
		},
	).File("/pull-requests.json").Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull requests: %w", err)
	}

	// Reprocess the result to simplify the JSON structure and only keeps
	// useful informations
	var pullRequests []*pullRequest
	if err := json.Unmarshal([]byte(result), &pullRequests); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pull requests: %w", err)
	}

	pullRequestsBytes, err := json.Marshal(pullRequests)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal pull requests: %w", err)
	}

	return dag.
		Directory().
		WithNewFile("/pull-requests.json", string(pullRequestsBytes)).
		File("/pull-requests.json"), nil
}
