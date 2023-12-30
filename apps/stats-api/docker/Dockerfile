ARG deps=mgarnier11/my-home:$TARGETARCH

FROM $deps as deps

FROM deps as build

WORKDIR /build

COPY --from=deps /build .
COPY ./apps/stats-api ./apps/stats-api
COPY --from=deps /apps/stat-api/node_modules ./apps/stats-api/node_modules

RUN npx nx build stats-api

# Use a multi-stage build to create a lean production image
FROM node:18-slim AS app

# Set the working directory
WORKDIR /app

COPY --from=build /build/apps/stats-api/dist ./dist

# Define the command to run the app
CMD ["node", "./dist/main.js"]
