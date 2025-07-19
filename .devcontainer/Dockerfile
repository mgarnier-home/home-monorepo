FROM mgarnier11/devcontainer-base:latest

COPY ./extensions.json /vscode/additionnal-extensions.json

RUN bash /vscode/install-extensions.sh "/vscode/additionnal-extensions.json"

COPY --chown=mgarnier:mgarnier ./01-container-config.zsh /setup/01-container-config.zsh
COPY --chown=mgarnier:mgarnier ./02-container-setup.zsh /setup/02-container-setup.zsh
RUN chmod +x /setup/01-container-config.zsh
RUN chmod +x /setup/02-container-setup.zsh

COPY --chown=mgarnier:mgarnier ./settings.json ./workspace.code-workspace /mnt/dev/.vscode/
