FROM build-go-app AS build


FROM registry.fedoraproject.org/fedora-minimal AS runtime

RUN microdnf -y install dnf-plugins-core && \
  dnf-3 config-manager --add-repo https://download.docker.com/linux/fedora/docker-ce.repo && \
  microdnf install -y docker-ce-cli docker-compose-plugin && microdnf clean all

RUN mkdir -p /config /config/entities/ /var/www/olivetin \
  && \
  microdnf install -y --nodocs --noplugins --setopt=keepcache=0 --setopt=install_weak_deps=0 \
  iputils \
  openssh-clients \
  jq \
  git \
  && microdnf clean all

WORKDIR /orchestrator

COPY --from=build /dist/* .

ENV COMPOSE_DIRECTORY="/compose"
ENV BINARIES_PATH="/orchestrator"
ENV SERVER_PORT="3000"

EXPOSE 3000

RUN mkdir -p /root/.ssh; \
  { echo "Host *"; \
  echo "    StrictHostKeyChecking no"; \
  echo "    IdentityFile /run/secrets/ssh_private_key"; } >> /root/.ssh/config; \
  chmod 600 /root/.ssh/config

ENTRYPOINT ["/orchestrator/orchestrator-api"]
