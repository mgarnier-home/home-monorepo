ARG deps=mgarnier11/my-home:$TARGETARCH

FROM $deps as deps

FROM deps as build

WORKDIR /build

COPY --from=deps /build .
COPY ./apps/node-proxy ./apps/node-proxy
COPY --from=deps /apps/node-proxy/node_modules ./apps/node-proxy/node_modules

RUN npx nx build node-proxy

RUN pnpm deploy --filter=node-proxy --prod /app

FROM node:20-alpine AS app

RUN apk add iputils

WORKDIR /app

COPY --from=build /app .

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget -q --spider --timeout=5 --tries=1 http://localhost:3000 || exit 1

EXPOSE 3000

ENV SERVER_PORT=3000
ENV CONFIG_FILE=/app/config.json

CMD ["node", "./dist/main.js"]
