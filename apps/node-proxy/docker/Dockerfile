FROM build AS build

FROM node:20-alpine AS app

RUN apk add iputils

WORKDIR /app

COPY --from=build /dist .

RUN ls -la

# HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
#   CMD wget -q --spider --timeout=5 --tries=1 http://localhost:3000 || exit 1

# ENV NODE_ENV="production"
# ENV SERVER_PORT="3000"

# EXPOSE 3000

# CMD ["node", "./main.js"]
