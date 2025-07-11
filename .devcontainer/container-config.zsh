#!/usr/bin/env zsh

echo "export ATHENA_HOST=ssh://mgarnier@100.64.98.100" >>~/.zshrc
echo "export ZEPHYR_HOST=ssh://mgarnier@100.64.98.97" >>~/.zshrc
echo "export APOLLON_HOST=ssh://mgarnier@100.64.98.99" >>~/.zshrc
echo "export ATLAS_HOST=ssh://mgarnier@100.64.98.98" >>~/.zshrc
echo "export COMPOSE_DIR=$docker_data_dir/zephyr/orchestrator" >>~/.zshrc
echo "export ENV_DIR=$docker_data_dir/zephyr/orchestrator" >>~/.zshrc
echo "source <(home-cli completion zsh)" >>~/.zshrc
echo "source <(task --completion zsh)" >>~/.zshrc
