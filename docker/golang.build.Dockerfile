FROM golang:1.22.3-alpine3.18

ARG APP=none
ARG VERSION=none

RUN apk add curl

WORKDIR /build

RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

COPY apps/$APP ./apps/$APP
COPY taskfile.yml .

RUN task ${APP}:build; \
  cp -r apps/$APP/dist /dist

