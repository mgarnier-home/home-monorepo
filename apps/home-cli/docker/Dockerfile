FROM build-go-app AS build


FROM jamesread/olivetin:latest AS runtime

USER root

RUN mkdir -p /usr/local/lib/docker/cli-plugins; \
  curl -SL https://github.com/docker/compose/releases/download/v2.29.2/docker-compose-linux-$(arch) -o /usr/local/lib/docker/cli-plugins/docker-compose; \
  chmod +x /usr/local/lib/docker/cli-plugins/docker-compose

# RUN microdnf install -y --nodocs --noplugins --setopt=keepcache=0 --setopt=install_weak_deps=0 \
#   ansible rsync

COPY --from=build /dist/home-cli /usr/bin/home-cli

ENV COMPOSE_DIR="/compose"
ENV OLIVETIN_CONFIG_DIR="/config"

RUN mkdir -p /root/.ssh; \
  { echo "Host *"; \
  echo "    StrictHostKeyChecking no"; \
  echo "    IdentityFile /run/secrets/ssh_private_key"; } >> /root/.ssh/config; \
  chmod 600 /root/.ssh/config

ENTRYPOINT ["/bin/sh", "-c", "/usr/bin/home-cli build-olivetin && /usr/bin/OliveTin"]
