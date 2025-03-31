# GitHub PR agent

A simple agent that can be used to simplify the review process on GitHub Pull requests.

## Usage

```
dagger

AGENT=$(. <repo> <gh-token>)

$AGENT | list-pull-requests

$AGENT | set-prnumber xxx | ask "What's the purpose of this PR?"
```

:bulb: I recommend to use the Dagger prompt mode to get a better experience. 

```
dagger

agent=$(. <repo> <gh-token>)

> # Enter prompt mode
# you can ask questions to the agent without using the shell syntax, it's much more natural
```

[Demo of the agent in action on the Dagger Discord](https://discord.com/channels/707636530424053791/1326978746703548416/1356239879033458739) 

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

The agent is using the [gh-repo-workspace](./gh-repo-workspace/) that exposes one tool:
- a tool list-pull-requests to get all the pull requests in the repository with the given filters.

The agent is using the [dagger-llm](https://docs.dagger.io/ai-agents) agent to interact with the conversation tool and the repository tool.