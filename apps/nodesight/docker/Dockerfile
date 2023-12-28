FROM node:18-alpine AS build
WORKDIR /build


COPY ./package.json ./pnpm-lock.yaml ./pnpm-workspace.yaml ./nx.json ./
COPY ./libs ./libs
COPY ./apps/nodesight ./apps/nodesight

RUN npm i -g pnpm

# Install dependencies
RUN pnpm install

RUN npx nx build nodesight

FROM nvidia/cuda:12.2.0-base-ubuntu20.04 AS deps-amd64

ENV NVIDIA_VISIBLE_DEVICES=all
ENV NVIDIA_DRIVER_CAPABILITIES="compute,video,utility"

RUN apt-get update && apt-get install -y ca-certificates curl gnupg && \
    mkdir -p /etc/apt/keyrings && \
    curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_18.x nodistro main" > /etc/apt/sources.list.d/nodesource.list && \
    apt-get update && apt-get install nodejs -y

RUN apt-get install pciutils dmidecode util-linux lm-sensors -y

FROM node:18-alpine AS deps-arm64

RUN apk add pciutils dmidecode util-linux lm-sensors

FROM deps-$TARGETARCH AS runtime

ARG TARGETARCH
ARG TARGETPLATFORM
ARG TARGETVARIANT

RUN echo "Target Arch : $TARGETARCH, Target Platform : $TARGETPLATFORM, Target Variant : $TARGETVARIANT"

# Set the working directory
WORKDIR /usr/app

# Copy compiled files from build stage
COPY --from=build /build/apps/nodesight/dist ./dist

# Define the command to run the app
CMD ["node", "./dist/main.js"]