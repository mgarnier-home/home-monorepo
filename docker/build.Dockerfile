FROM node:20-alpine

ARG APP=none
ARG APP_VERSION=none

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
# Enable corepack, that allows us to use pnpm as a drop-in replacement for npm
RUN corepack enable

WORKDIR /build

COPY apps/$APP ./apps/$APP
COPY libs ./libs
COPY pnpm-workspace.yaml package.json pnpm-lock.yaml tsconfig.json default.webpack.config.ts ./


RUN echo $APP

RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

RUN pnpm run --filter=$APP build

RUN cp -r apps/$APP/dist /dist

RUN echo $APP_VERSION > /dist/version.txt