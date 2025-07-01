FROM mgarnier11/devcontainer-base:latest

COPY ./extensions.json /vscode/additionnal-extensions.json

RUN bash /vscode/install-extensions.sh "/vscode/additionnal-extensions.json"

COPY --chown=mgarnier:mgarnier ./container-config.fish /home/mgarnier/.container-config.fish
COPY --chown=mgarnier:mgarnier ./container-setup.fish /setup/container-setup.fish
RUN chmod +x /setup/container-setup.fish

COPY --chown=mgarnier:mgarnier ./settings.json ./workspace.code-workspace /mnt/dev/.vscode/
