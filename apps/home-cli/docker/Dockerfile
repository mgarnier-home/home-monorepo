FROM build AS build


FROM jamesread/olivetin:latest AS runtime

USER root

# RUN microdnf install -y --nodocs --noplugins --setopt=keepcache=0 --setopt=install_weak_deps=0 \
#   ansible rsync

COPY --from=build /dist/home-cli /usr/bin/home-cli

ENV COMPOSE_DIR="/compose"
ENV OLIVETIN_CONFIG_DIR="/config"

ENTRYPOINT ["/bin/sh", "-c", "/usr/bin/home-cli build-olivetin && /usr/bin/OliveTin"] 