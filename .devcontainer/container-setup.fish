#!/usr/bin/env fish

function setupHomeMonorepo
    set home_monorepo_dir /mnt/dev/home-monorepo
    set docker_data_dir /mnt/docker-data

    if not test -d $home_monorepo_dir
        echo "Cloning home-monorepo"
        git clone git@github.com:mgarnier-home/home-monorepo.git $home_monorepo_dir
    end

    cd $home_monorepo_dir

    git config --global --add safe.directory $home_monorepo_dir

    # Install pnpm and install dependencies
    npm install pnpm -g
    pnpm approve-builds
    pnpm install

    task home-cli:build

    if not test -L ~/.local/home-cli
        echo "Creating symlink for home-cli"
        ln -s $home_monorepo_dir/apps/home-cli/dist/home-cli ~/.local/home-cli
    else
        echo "Symlink for home-cli already exists"
    end
end

setupHomeMonorepo

echo "Setup done!"
