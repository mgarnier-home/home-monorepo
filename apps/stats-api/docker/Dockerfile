# Use the official lightweight Node.js 18.11.0 image.
# https://hub.docker.com/_/node
FROM node:18.11.0-alpine AS build

# Set the working directory
WORKDIR /usr/app

RUN npm i -g pnpm

# Install dependencies
COPY ./package.json ./pnpm-lock.yaml ./
RUN pnpm install

# Copy local code to the container image
COPY ./src ./src
COPY ./tsconfig.json ./tsconfig.json

# Build the production application
RUN pnpm run compile

# Remove source files
RUN rm -rf ./src

# Remove development dependencies
RUN pnpm prune --prod

# Use a multi-stage build to create a lean production image
FROM node:18.11.0-alpine AS runtime

# Set the working directory
WORKDIR /usr/app

# Copy package.json and package-lock.json for necessary production dependencies
COPY --from=build /usr/app/package.json /usr/app/pnpm-lock.yaml ./

# Copy node modules from build stage
COPY --from=build /usr/app/node_modules ./node_modules

# Copy compiled files from build stage
COPY --from=build /usr/app/dist ./dist

# Define the command to run the app
CMD ["node", "./dist/index.js"]
