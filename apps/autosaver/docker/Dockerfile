ARG deps=mgarnier11/my-home:$TARGETARCH

FROM $deps as deps

FROM deps as build

WORKDIR /build

COPY --from=deps /build .
COPY ./apps/autosaver ./apps/autosaver
COPY --from=deps /apps/autosaver/node_modules ./apps/autosaver/node_modules

RUN npx nx build autosaver

RUN pnpm deploy --filter=autosaver --prod /app

# Use a multi-stage build to create a lean production image
FROM node:18.11.0-buster AS runtime

RUN apt update && apt install p7zip-full rsync cifs-utils -y

# Set the working directory
WORKDIR /app

COPY --from=build /app .

# Define the command to run the app
CMD ["node", "./dist/main.js"]
