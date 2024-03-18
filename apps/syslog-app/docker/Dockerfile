ARG deps=mgarnier11/my-home:$TARGETARCH

FROM $deps as deps

FROM deps as build

WORKDIR /build

COPY --from=deps /build .
COPY ./apps/syslog-app ./apps/syslog-app
COPY --from=deps /apps/syslog-app/node_modules ./apps/syslog-app/node_modules

RUN npx nx build syslog-app

FROM node:18-alpine AS app

ENV NODE_ENV=production

WORKDIR /app

COPY --from=build /build/apps/syslog-app/server-dist ./server-dist
COPY --from=build /build/apps/syslog-app/app-dist ./app-dist

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget -q --spider --timeout=5 --tries=1 http://localhost:3000 || exit 1

CMD ["node", "./server-dist/main.js"]