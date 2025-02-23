#!/usr/bin/zsh
. ~/.zshrc

set -euo pipefail

cp /run/secrets/SSH_PRIVATE_KEY ~/.ssh/id_rsa
chmod 600 ~/.ssh/id_rsa

git config --global core.autocrlf false
git config --global user.email "mgarnier11@gmail.com"
git config --global user.name "Mathieu GARNIER"

setupHomeMonorepo() {
  home_monorepo_dir=/mnt/dev/home-monorepo
  docker_data_dir=/mnt/docker-data

  if [ ! -d $home_monorepo_dir ]; then
    echo "Cloning home-monorepo"
    git clone git@github.com:mgarnier-home/home-monorepo.git $home_monorepo_dir
  fi

  cd $home_monorepo_dir

  git config --global --add safe.directory $home_monorepo_dir

  # Install pnpm and install dependencies
  npm install pnpm -g
  pnpm install

  # Build and install home-cli
  export COMPOSE_DIR="$docker_data_dir/zephyr/orchestrator"
  export ENV_DIR="$docker_data_dir/zephyr/orchestrator"

  task home-cli:build

  ln -s $home_monorepo_dir/apps/home-cli/dist/home-cli ~/.local/home-cli

  echo "export ATHENA_HOST=ssh://mgarnier@100.64.98.100" >>~/.zshrc
  echo "export ZEPHYR_HOST=ssh://mgarnier@100.64.98.97" >>~/.zshrc
  echo "export APOLLON_HOST=ssh://mgarnier@100.64.98.99" >>~/.zshrc
  echo "export COMPOSE_DIR=$docker_data_dir/zephyr/orchestrator" >>~/.zshrc
  echo "export ENV_DIR=$docker_data_dir/zephyr/orchestrator" >>~/.zshrc
  echo "source <(home-cli completion zsh)" >>~/.zshrc
}

setupGhActions() {
  gh_actions_dir=/mnt/dev/gh-actions

  if [ ! -d $gh_actions_dir ]; then
    echo "Cloning gh-actions"
    git clone git@github.com:mgarnier11/gh-actions.git $gh_actions_dir
  fi

  git config --global --add safe.directory $gh_actions_dir

  cd $gh_actions_dir

  git checkout dev
}

setupBlindtestGen() {
  blindtest_gen_dir=/mnt/dev/blindtest-gen

  if [ ! -d $blindtest_gen_dir ]; then
    echo "Cloning blindtest-gen"
    git clone git@github.com:mgarnier11/blindtest-gen.git $blindtest_gen_dir
  fi

  git config --global --add safe.directory $blindtest_gen_dir
}


setupVsCodeFolder() {
  vscode_dir=/mnt/dev/.vscode

  mkdir -p $vscode_dir

  mv /setup/settings.json $vscode_dir/settings.json
  mv /setup/workspace.code-workspace $vscode_dir/workspace.code-workspace
}

setupHomeMonorepo
setupGhActions
setupBlindtestGen
setupVsCodeFolder

echo "Setup done!"
