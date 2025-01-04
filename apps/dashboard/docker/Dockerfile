FROM build-node-app AS build

FROM node:20-alpine AS runtime

WORKDIR /app

COPY --from=build /dist/server ./server
COPY --from=build /dist/app/browser ./app


HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget -q --spider --timeout=5 --tries=1 http://localhost:3000 || exit 1

ENV NODE_ENV="production"
ENV SERVER_PORT="3000"
ENV APP_DIST_PATH="/app/app"

EXPOSE 3000

CMD ["node", "./server/main.js"]
