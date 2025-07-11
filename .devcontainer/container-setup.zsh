#!/usr/bin/env zsh

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

    # check if the symlink  ~/.local/home-cli exists and create a symlink if it does not
    if [ ! -d ~/.local/home-cli ]; then
        echo "Creating symlink for home-cli"
        mkdir -p ~/.local
        ln -s $home_monorepo_dir/apps/home-cli/dist/home-cli ~/.local/home-cli
    else
        echo "Symlink for home-cli already exists"
    fi

}

setupHomeMonorepo

echo "Setup done!"
