#!/usr/bin/env zsh


mkdir -p /home/mgarnier/.local
curl -o "/home/mgarnier/.local/orchestrator-cli" "$API_ORCHESTRATOR_URL/cli?os=linux&arch=amd64"
chmod +x /home/mgarnier/.local/orchestrator-cli

echo "export API_ORCHESTRATOR_URL=$API_ORCHESTRATOR_URL" >>~/.zshrc
echo "source <(orchestrator-cli completion zsh)" >>~/.zshrc
echo "source <(task --completion zsh)" >>~/.zshrc
