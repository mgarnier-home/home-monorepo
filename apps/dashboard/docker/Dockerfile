FROM node:18.11.0-alpine AS build
WORKDIR /build

COPY ./package.json ./pnpm-lock.yaml ./pnpm-workspace.yaml ./nx.json ./
COPY ./libs ./libs
COPY ./apps/dashboard ./apps/dashboard

# Install pnpm
RUN npm i -g pnpm

# Install dependencies
RUN pnpm install

# Build the app
RUN npx nx build dashboard

FROM node:18.11.0-alpine AS app

ENV NODE_ENV=production

WORKDIR /app

COPY --from=build /build/apps/dashboard/server-dist ./server-dist
COPY --from=build /build/apps/dashboard/app-dist ./app-dist

CMD ["node", "./server-dist/main.js"]
