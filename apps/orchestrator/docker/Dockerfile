FROM build-go-app AS build


FROM debian:bookworm-slim AS runtime

RUN apt-get update && \
  apt-get install -y \
  jq \
  curl \
  gnupg \
  lsb-release \
  openssh-client \
  ca-certificates \
  apt-transport-https && \
  curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg && \
  echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian bookworm stable" \
  > /etc/apt/sources.list.d/docker.list && \
  apt-get update && \
  apt-get install -y docker-ce-cli docker-compose-plugin && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*


WORKDIR /orchestrator

COPY --from=build /dist/* .

ENV COMPOSE_DIRECTORY="/compose"
ENV BINARIES_PATH="/orchestrator"

RUN mkdir -p /root/.ssh; \
  { echo "Host *"; \
  echo "    StrictHostKeyChecking no"; \
  echo "    IdentityFile /run/secrets/ssh_private_key"; } >> /root/.ssh/config; \
  chmod 600 /root/.ssh/config

ENTRYPOINT ["/orchestrator/orchestrator-api"]
