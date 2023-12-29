FROM mgarnier11/my-home-deps as build

WORKDIR /build

COPY --from=mgarnier11/my-home-deps /build .

RUN npx nx build stats-api

# Use a multi-stage build to create a lean production image
FROM node:18-slim AS runtime

# Set the working directory
WORKDIR /usr/app

COPY --from=build /build/apps/stats-api/dist ./dist

# Define the command to run the app
CMD ["node", "./dist/main.js"]
