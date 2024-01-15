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

CMD ["node", "./server-dist/main.js"]