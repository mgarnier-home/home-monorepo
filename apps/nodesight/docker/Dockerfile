FROM build AS build

FROM nvidia/cuda:12.2.0-base-ubuntu20.04 AS deps-amd64-nvidia

ENV NVIDIA_VISIBLE_DEVICES=all
ENV NVIDIA_DRIVER_CAPABILITIES="compute,video,utility"

RUN apt-get update && apt-get install -y ca-certificates curl gnupg && \
  mkdir -p /etc/apt/keyrings && \
  curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
  echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_20.x nodistro main" > /etc/apt/sources.list.d/nodesource.list && \
  apt-get update && apt-get install nodejs -y

RUN apt-get install pciutils dmidecode util-linux lm-sensors -y

FROM node:20-alpine AS deps-amd64

RUN apk add pciutils dmidecode util-linux lm-sensors

FROM node:20-alpine AS deps-arm64

RUN apk add pciutils dmidecode util-linux lm-sensors

FROM deps-$TARGETARCH AS app

ARG TARGETARCH
ARG TARGETPLATFORM
ARG TARGETVARIANT

RUN echo "Target Arch : $TARGETARCH, Target Platform : $TARGETPLATFORM, Target Variant : $TARGETVARIANT"

# Set the working directory
WORKDIR /app

# Copy compiled files from build stage
COPY --from=build /dist .

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget -q --spider --timeout=5 --tries=1 http://localhost:3000 || exit 1

ENV NODE_ENV="production"
ENV SERVER_PORT="3000"
ENV UPDATE_INTERVAL="15000"
ENV ENABLE_STATS_API="false"
ENV DISABLE_CPU_TEMPS="false"

EXPOSE 3000

# Define the command to run the app
CMD ["node", "./main.js"]