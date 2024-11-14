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

# setupHomeContainer() {
#   home_container_dir=/mnt/dev/home-container

#   if [ ! -d $home_container_dir ]; then
#     git clone git@github.com:mgarnier-home/home-container.git $home_container_dir
#   fi

#   git config --global --add safe.directory $home_container_dir
# }

# setupHomeConfig() {
#   home_config_dir=/mnt/dev/home-config

#   if [ ! -d $home_config_dir ]; then
#     git clone git@github.com:mgarnier-home/home-config.git $home_config_dir
#   fi

#   git config --global --add safe.directory $home_config_dir
# }

# setupCoderTemplates() {
#   CODER_TEMPLATES_DIR=/mnt/dev/coder-templates

#   if [ ! -d $CODER_TEMPLATES_DIR ]; then
#     git clone git@github.com:mgarnier-home/coder-templates.git $CODER_TEMPLATES_DIR
#   fi

#   git config --global --add safe.directory $CODER_TEMPLATES_DIR

#   cd $CODER_TEMPLATES_DIR

#   coder login --url https://coder.int.mgarnier11.fr --token ${CODER_SESSION_TOKEN}
# }

setupGhActions() {
  gh_actions_dir=/mnt/dev/gh-actions

  if [ ! -d $gh_actions_dir ]; then
    git clone git@github.com:mgarnier11/gh-actions.git $gh_actions_dir
  fi

  git config --global --add safe.directory $gh_actions_dir

  cd $gh_actions_dir

  git checkout dev
}

# setupTerraformModules() {
#   TERRAFORM_MODULES_DIR=/mnt/dev/terraform-modules

#   if [ ! -d $TERRAFORM_MODULES_DIR ]; then
#     git clone git@github.com:mgarnier-home/terraform-modules.git $TERRAFORM_MODULES_DIR
#   fi

#   git config --global --add safe.directory $TERRAFORM_MODULES_DIR
# }

# setupMyTraefik() {
#   MY_TRAEFIK_DIR=/mnt/dev/my-traefik

#   if [ ! -d $MY_TRAEFIK_DIR ]; then
#     git clone git@github.com:mgarnier-home/my-traefik.git $MY_TRAEFIK_DIR
#   fi

#   git config --global --add safe.directory $MY_TRAEFIK_DIR
# }

setupHomeMonorepo
setupGhActions

bash /setup/get-workspace-file.sh "$(tr '\n' ' ' </setup/workspace.json)"

echo "Setup done!"
