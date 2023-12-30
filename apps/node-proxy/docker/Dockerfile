ARG deps=mgarnier11/my-home:$TARGETARCH

FROM $deps as deps

FROM deps as build

WORKDIR /build

COPY --from=deps /build .
COPY ./apps/node-proxy ./apps/node-proxy
COPY --from=deps /apps/node-proxy/node_modules ./apps/node-proxy/node_modules

RUN pnpm install --shamefully-hoist --force

RUN npx nx build node-proxy

RUN pnpm prune --prod

# Use a multi-stage build to create a lean production image
FROM node:18-slim AS app

# Set the working directory
WORKDIR /app

COPY --from=build /build/node_modules ./node_modules
COPY --from=build /build/apps/node-proxy/dist ./dist

# Define the command to run the app
CMD ["node", "./dist/main.js"]


# FROM node:18.11.0-alpine AS build

# WORKDIR /usr/app

# RUN npm i -g pnpm

# # Install dependencies
# COPY ./package.json ./
# RUN pnpm install --shamefully-hoist

# # Copy local code to the container image
# COPY ./src ./src
# COPY ./tsconfig.json ./tsconfig.json

# # Compile the code
# RUN pnpm run compile

# RUN rm -rf ./src

# RUN pnpm prune --prod

# FROM node:18.11.0-alpine AS app

# RUN apk add iputils

# WORKDIR /usr/app

# COPY --from=build /usr/app/ .

# CMD node ./dist/index.js
