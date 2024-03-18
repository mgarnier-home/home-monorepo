ARG deps=mgarnier11/my-home:$TARGETARCH

FROM $deps as deps

FROM deps as build

WORKDIR /build

COPY --from=deps /build .
COPY ./apps/traefik-conf ./apps/traefik-conf
COPY --from=deps /apps/traefik-conf/node_modules ./apps/traefik-conf/node_modules

RUN npx nx build traefik-conf

# Use a multi-stage build to create a lean production image
FROM node:20-alpine AS app

# Set the working directory
WORKDIR /app

COPY --from=build /build/apps/traefik-conf/dist ./dist

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget -q --spider --timeout=5 --tries=1 http://localhost:3000 || exit 1

# Define the command to run the app
CMD ["node", "./dist/main.js"]
