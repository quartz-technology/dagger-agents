package main

import (
	"context"
	"dagger/gh-pr-workspace/internal/dagger"
	"encoding/json"
	"fmt"
)

type PRConversation struct {
	Comments []*PRComment
	Reviews  []*PRReview
}

type PRComment struct {
	Author     string `json:"author"`
	Body       string `json:"body"`
	Timestamps string `json:"created_at"`
}

type PRReview struct {
	Author     string `json:"author"`
	Body       string `json:"body"`
	Timestamps string `json:"created_at"`
	Diff       string `json:"diff"`
	File       string `json:"filename"`
}

// Custom UnmarshalJSON function for PRComment
func (c *PRComment) UnmarshalJSON(data []byte) error {
	// Define a temporary struct to match the full JSON structure
	var temp struct {
		User struct {
			Login string `json:"login"`
		} `json:"user"`
		Body      string `json:"body"`
		Timestamp string `json:"created_at"`
	}

	// Parse JSON into the temporary struct
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Extract only the required fields
	c.Author = temp.User.Login
	c.Body = temp.Body
	c.Timestamps = temp.Timestamp

	return nil
}

func (c *PRReview) UnmarshalJSON(data []byte) error {
	// Define a temporary struct to match the full JSON structure
	var temp struct {
		User struct {
			Login string `json:"login"`
		} `json:"user"`
		Body      string `json:"body"`
		Timestamp string `json:"created_at"`
		Diff      string `json:"diff_hunk"`
		File      string `json:"path"`
	}

	// Parse JSON into the temporary struct
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Extract only the required fields
	c.Author = temp.User.Login
	c.Body = temp.Body
	c.Timestamps = temp.Timestamp
	c.Diff = temp.Diff
	c.File = temp.File

	return nil
}

// Conversation returns a file containing all the messages sent in the PR in a structured format.
func (g *GhPrWorkspace) Conversation(ctx context.Context) (*dagger.File, error) {
	prComments, err := g.comments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR comments: %w", err)
	}

	prReviews, err := g.reviews(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR reviews: %w", err)
	}

	prConversation := PRConversation{
		Comments: prComments,
		Reviews:  prReviews,
	}

	prConversationsBytes, err := json.Marshal(prConversation)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal PR messages: %w", err)
	}

	return dag.
		Directory().
		WithNewFile("/pr-messages.json", string(prConversationsBytes)).
		File("/pr-messages.json"), nil
}

func (g *GhPrWorkspace) comments(ctx context.Context) ([]*PRComment, error) {
	result, err := g.Container.WithExec(
		[]string{"gh", "api", fmt.Sprintf("repos/%s/issues/%d/comments", g.RepositoryURL, g.PullRequestNumber)},
		dagger.ContainerWithExecOpts{
			RedirectStdout: "/comments.json",
		},
	).File("/comments.json").Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %w", err)
	}

	var comments []*PRComment
	if err := json.Unmarshal([]byte(result), &comments); err != nil {
		return nil, fmt.Errorf("failed to unmarshal comments: %w", err)
	}

	return comments, nil
}

func (g *GhPrWorkspace) reviews(ctx context.Context) ([]*PRReview, error) {
	result, err := g.Container.WithExec(
		[]string{"gh", "api", fmt.Sprintf("repos/%s/pulls/%d/comments", g.RepositoryURL, g.PullRequestNumber)},
		dagger.ContainerWithExecOpts{
			RedirectStdout: "/reviews.json",
		},
	).File("/reviews.json").Contents(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get review comments: %w", err)
	}

	var reviews []*PRReview
	if err := json.Unmarshal([]byte(result), &reviews); err != nil {
		return nil, fmt.Errorf("failed to unmarshal reviews: %w", err)
	}

	return reviews, nil
}
