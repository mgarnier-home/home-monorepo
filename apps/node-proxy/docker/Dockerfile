ARG deps=mgarnier11/my-home:$TARGETARCH

FROM $deps as deps

FROM deps as build

WORKDIR /build

COPY --from=deps /build .

COPY ./apps/node-proxy ./apps/node-proxy

WORKDIR /build/apps/node-proxy

RUN pnpm install --shamefully-hoist

RUN pnpm run compile



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
