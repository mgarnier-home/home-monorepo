FROM node:20-bookworm

ARG APP=none
ARG VERSION=none

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
# Enable corepack, that allows us to use pnpm as a drop-in replacement for npm
RUN corepack enable

# RUN apt-get update && apt-get install build-essential -y

WORKDIR /build

COPY apps/$APP ./apps/$APP
COPY libs ./libs
COPY pnpm-workspace.yaml package.json pnpm-lock.yaml tsconfig.json default.webpack.config.ts .pnpmfile.cjs ./

RUN pnpm install --frozen-lockfile

RUN pnpm run --filter=$APP build

RUN cp -r apps/$APP/dist /dist

RUN echo $VERSION > /dist/version.txt
