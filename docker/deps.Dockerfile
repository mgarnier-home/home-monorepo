FROM node:18-slim

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
# Enable corepack, that allows us to use pnpm as a drop-in replacement for npm
RUN corepack enable

WORKDIR /build

COPY . .

RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile --prod
RUN mv ./apps /apps
