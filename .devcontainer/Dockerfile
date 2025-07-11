FROM mgarnier11/devcontainer-base:latest

COPY ./extensions.json /vscode/additionnal-extensions.json

RUN bash /vscode/install-extensions.sh "/vscode/additionnal-extensions.json"

COPY --chown=mgarnier:mgarnier ./container-config.zsh /setup/container-config.zsh
COPY --chown=mgarnier:mgarnier ./container-setup.zsh /setup/container-setup.zsh
RUN chmod +x /setup/container-setup.zsh
RUN chmod +x /setup/container-config.zsh

COPY --chown=mgarnier:mgarnier ./settings.json ./workspace.code-workspace /mnt/dev/.vscode/
