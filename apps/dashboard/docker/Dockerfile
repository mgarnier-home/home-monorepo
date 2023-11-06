FROM node:18.11.0-alpine AS build
WORKDIR /build

COPY ./package.json ./pnpm-lock.yaml ./pnpm-workspace.yaml ./
COPY ./libs ./libs
COPY ./apps/dashboard ./apps/dashboard

# Install pnpm
RUN npm i -g pnpm

WORKDIR /build/apps/dashboard

# Install dependencies
RUN pnpm install

# Build the app
RUN pnpm run build

# Prune dev dependencies
RUN pnpm prune --prod

RUN rm -rf ./server ./shared ./app

# Install production dependencies
ENV NODE_ENV=production

CMD ["node", "./server-dist/server/main.js"]
