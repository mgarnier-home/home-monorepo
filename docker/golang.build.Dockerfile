FROM golang:1.24.5-alpine

ARG APP=none
ARG VERSION=none

RUN apk add curl

WORKDIR /build

RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin

COPY apps/$APP ./apps/$APP
COPY libs/go ./libs/go
COPY taskfile.yml .

RUN \ 
  task ${APP}:build && \
  cp -r apps/$APP/dist /dist && \
  echo $VERSION > /dist/version.txt

