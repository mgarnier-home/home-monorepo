FROM mgarnier11/my-home-deps as build

WORKDIR /build

COPY --from=mgarnier11/my-home-deps /build .

# RUN pnpm run --filter dashboard build

WORKDIR /build/apps/dashboard

RUN pnpm run esbuild:app -- build
RUN pnpm run build:app:tailwind
RUN pnpm run lint:server
RUN pnpm run esbuild:server -- build

FROM node:18-slim AS dashboard-app

ENV NODE_ENV=production

WORKDIR /app

COPY --from=build /build/apps/dashboard/server-dist ./server-dist
COPY --from=build /build/apps/dashboard/app-dist ./app-dist

CMD ["node", "./server-dist/main.js"]
