#!/usr/bin/env zsh

. ~/.zshrc


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
}

setupHomeMonorepo

echo "Setup done!"
