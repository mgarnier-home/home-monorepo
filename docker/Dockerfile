FROM node:20-alpine

ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
# Enable corepack, that allows us to use pnpm as a drop-in replacement for npm
RUN corepack enable

WORKDIR /build

COPY . .

# Install dependencies for all workspaces
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install --frozen-lockfile

# Move apps to a separate directory, so that we can copy the /build directory without having to copy possibly not up to date apps
RUN mv ./apps /apps
