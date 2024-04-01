ARG deps=mgarnier11/my-home:$TARGETARCH

FROM $deps as deps

FROM deps as build

WORKDIR /build

COPY --from=deps /build .
COPY ./apps/orchestrator ./apps/orchestrator
COPY --from=deps /apps/orchestrator/node_modules ./apps/orchestrator/node_modules

RUN npx nx build orchestrator

# Use a multi-stage build to create a lean production image
FROM node:20-alpine AS app

RUN apk update 
RUN apk add \
  docker \
  docker-compose \
  openssh-client \
  ca-certificates

RUN mkdir -p /root/.ssh
RUN echo "StrictHostKeyChecking no" >> /root/.ssh/config

# Set the working directory
WORKDIR /app

COPY --from=build /build/apps/orchestrator/dist ./dist

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget -q --spider --timeout=5 --tries=1 http://localhost:3000 || exit 1

# Define the command to run the app
CMD ["node", "./dist/main.js"]
