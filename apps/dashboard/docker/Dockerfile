FROM node:18.11.0-alpine AS build
WORKDIR /usr/app

# Install pnpm
RUN npm i -g pnpm

# Copy required files
COPY ./package.json ./pnpm-lock.yaml ./tailwind.config.ts ./postcss.config.js ./tsconfig.json ./vite.config.ts ./
COPY ./app ./app
COPY ./shared ./shared

# Install dependencies
RUN pnpm install

# Build the app
RUN pnpm run build

FROM node:18.11.0-alpine AS runtime
WORKDIR /usr/app

# Copy node modules and built app from build stage
COPY --from=build /usr/app/node_modules ./node_modules
COPY --from=build /usr/app/app-dist ./app-dist

# Copy required runtime files
COPY ./package.json ./pnpm-lock.yaml ./vite.config.ts ./tsconfig.json ./
COPY ./server ./server
COPY ./shared ./shared

ENV NODE_ENV=production

CMD ["npm", "run", "start"]
