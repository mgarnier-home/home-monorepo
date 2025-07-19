#!/usr/bin/env zsh


mkdir -p /home/mgarnier/.local
curl -o "/home/mgarnier/.local/orchestrator-cli" "$ORCHESTRATOR_API_URL/cli?os=linux&arch=amd64"
chmod +x /home/mgarnier/.local/orchestrator-cli

echo "export ORCHESTRATOR_API_URL=$ORCHESTRATOR_API_URL" >>~/.zshrc
echo "source <(orchestrator-cli completion zsh)" >>~/.zshrc
echo "source <(task --completion zsh)" >>~/.zshrc
