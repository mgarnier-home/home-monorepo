ARG deps=mgarnier11/my-home:$TARGETARCH

FROM $deps as deps

FROM deps as build

WORKDIR /build

COPY --from=deps /build .
COPY ./apps/dashboard ./apps/dashboard
COPY --from=deps /apps/dashboard/node_modules ./apps/dashboard/node_modules

RUN npx nx build dashboard

FROM node:18-alpine AS app

ENV NODE_ENV=production

WORKDIR /app

COPY --from=build /build/apps/dashboard/server-dist ./server-dist
COPY --from=build /build/apps/dashboard/app-dist ./app-dist

CMD ["node", "./server-dist/main.js"]