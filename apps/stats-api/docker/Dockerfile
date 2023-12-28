# Use the official lightweight Node.js 18.11.0 image.
# https://hub.docker.com/_/node
FROM node:18.11.0-alpine AS build

WORKDIR /build

COPY ./package.json ./pnpm-lock.yaml ./pnpm-workspace.yaml ./nx.json ./
COPY ./libs ./libs
COPY ./apps/stats-api ./apps/stats-api

RUN npm i -g pnpm

# Install dependencies
RUN pnpm install

RUN npx nx build stats-api

# Use a multi-stage build to create a lean production image
FROM node:18.11.0-alpine AS runtime

# Set the working directory
WORKDIR /usr/app

COPY --from=build /build/apps/stats-api/dist ./dist

# Define the command to run the app
CMD ["node", "./dist/main.js"]
