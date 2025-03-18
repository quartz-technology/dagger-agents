# GitHub PR agent

A simple agent that can be used to simplify the review process on GitHub Pull requests.

## Usage

```
dagger

AGENT=$(. <repo> <gh-token>)

$AGENT | set-prnumber xxx | ask "What's the purpose of this PR?"
```

## Examples of prompts

```
$AGENT | set-prnumber xxx | ask "What files have changed in this PR?"

$AGENT | set-prnumber xxx | ask "Can you catch all mistakes that should be fixed on the PR?"

$AGENT | set-prnumber xxx | ask "Can you verify that the documentation is written in correct english? Can you suggest improvements or corrections?"
```

## Technical details.

The agent is using the [gh-pr-workspace](./gh-pr-workspace/) that exposes two tools:
- a tool conversation to get all the messages sent in the PR including reviews and comments.
- a tool repository to get the differences between the PR and the origin or read files contents.

The agent is using the [dagger-llm](https://docs.dagger.io/ai-agents) agent to interact with the conversation tool and the repository tool.