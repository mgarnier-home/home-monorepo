FROM node:18-slim AS dashboard-app

ENV NODE_ENV=production

WORKDIR /app

COPY --from=mgarnier11/dashboard:dashboard-build /build/apps/dashboard/server-dist ./server-dist
COPY --from=mgarnier11/dashboard:dashboard-build /build/apps/dashboard/app-dist ./app-dist

CMD ["node", "./server-dist/main.js"]