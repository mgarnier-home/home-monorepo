cp /run/secrets/SSH_PRIVATE_KEY ~/.ssh/id_rsa
chmod 600 ~/.ssh/id_rsa

git config --global core.autocrlf false
git config --global user.email "mgarnier11@gmail.com"
git config --global user.name "Mathieu GARNIER"

# install node
set -q NODE_VERSION; or set NODE_VERSION 20

nvm install $NODE_VERSION
nvm use $NODE_VERSION


# install go
set -q GO_VERSION; or set GO_VERSION 1.24.4
curl -o go.tar.gz https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go.tar.gz
rm go.tar.gz
set -U fish_user_paths /usr/local/go/bin $fish_user_paths

function setupVsCodeFolder
    set vscode_dir /mnt/dev/.vscode

    mkdir -p $vscode_dir

    mv /setup/settings.json $vscode_dir/settings.json
    mv /setup/workspace.code-workspace $vscode_dir/workspace.code-workspace
end

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
    pnpm install

    task home-cli:build

    if not test -L ~/.local/home-cli
        echo "Creating symlink for home-cli"
        ln -s $home_monorepo_dir/apps/home-cli/dist/home-cli ~/.local/home-cli
    else
        echo "Symlink for home-cli already exists"
    end
end

setupVsCodeFolder
setupHomeMonorepo


echo "Setup done!"
