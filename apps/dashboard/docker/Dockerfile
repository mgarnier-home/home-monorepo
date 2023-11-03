# syntax = edrevo/dockerfile-plus
ARG APP=dashboard
FROM node:18.11.0-alpine AS build
WORKDIR /build

INCLUDE+ Dockerfile.base

# Install pnpm
RUN npm i -g pnpm

WORKDIR /build/apps/dashboard

# Install dependencies
RUN pnpm install

# Build the app
RUN pnpm run build

FROM node:18.11.0-alpine AS runtime
WORKDIR /usr/app

# Copy node modules and built app from build stage
COPY --from=build /build/apps/$APP/node_modules ./node_modules
COPY --from=build /build/apps/$APP/app-dist ./app-dist

# Copy required runtime files
COPY ./apps/$APP/package.json ./apps/$APP/vite.config.ts ./apps/$APP/tsconfig.json ./
COPY ./apps/$APP/server ./server
COPY ./apps/$APP/shared ./shared

ENV NODE_ENV=production

CMD ["npm", "run", "start"]
