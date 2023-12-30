FROM mgarnier11/my-home:$TARGETARCH as dashboard-build

WORKDIR /build

COPY --from=mgarnier11/my-home /build .
COPY ./apps/dashboard ./apps/dashboard
COPY --from=mgarnier11/my-home /apps/dashboard/node_modules ./apps/dashboard/node_modules

RUN npx nx build dashboard

FROM node:18-slim AS dashboard-app

ENV NODE_ENV=production

WORKDIR /app

COPY --from=dashboard-build /build/apps/dashboard/server-dist ./server-dist
COPY --from=dashboard-build /build/apps/dashboard/app-dist ./app-dist

CMD ["node", "./server-dist/main.js"]