FROM node:18-slim AS dashboard-app

ENV NODE_ENV=production

WORKDIR /app

COPY ./apps/dashboard/server-dist ./server-dist
COPY ./apps/dashboard/app-dist ./app-dist

CMD ["node", "./server-dist/main.js"]