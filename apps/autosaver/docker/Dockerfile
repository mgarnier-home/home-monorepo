FROM build AS build

FROM node:20-buster AS runtime

RUN apt update && apt install p7zip-full rsync cifs-utils gpg -y

RUN mkdir -p /root/.gnupg && touch /root/.gnupg/pubring.kbx && chmod 700 /root/.gnupg

# Set the working directory
WORKDIR /app

COPY --from=build /dist .

RUN touch config.yml

RUN ls -la

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget -q --spider --timeout=5 --tries=1 http://localhost:3000 || exit 1

ENV NODE_ENV="production"
ENV SERVER_PORT="3000"
ENV CRON_SCHEDULE="0 0 * * *"
ENV ARCHIVE_API_TYPE="zip"

EXPOSE 3000

# Define the command to run the app
CMD ["node", "./main.js"]
