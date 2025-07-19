#!/usr/bin/env zsh

curl -o ~/.local/orchestrator-cli $ORCHESTRATOR_API_URL/cli

echo "source <(orchestrator-cli completion zsh)" >>~/.zshrc
echo "source <(task --completion zsh)" >>~/.zshrc
