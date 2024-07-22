FROM build AS build

# Use a multi-stage build to create a lean production image
FROM node:20-alpine AS app

# Set the working directory
WORKDIR /app

COPY --from=build /dist .

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget -q --spider --timeout=5 --tries=1 http://localhost:3000 || exit 1

ENV NODE_ENV="production"
ENV SERVER_PORT="3000"

EXPOSE 3000

# Define the command to run the app
CMD ["node", "./dist/main.js"]
