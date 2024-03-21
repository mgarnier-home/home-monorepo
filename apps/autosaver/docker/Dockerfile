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
FROM node:20-buster AS runtime

RUN apt update && apt install p7zip-full rsync cifs-utils gpg -y

RUN mkdir -p /root/.gnupg && touch /root/.gnupg/pubring.kbx

# Set the working directory
WORKDIR /app

COPY --from=build /app .

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget -q --spider --timeout=5 --tries=1 http://localhost:3000 || exit 1

# Define the command to run the app
CMD ["node", "./dist/main.js"]
