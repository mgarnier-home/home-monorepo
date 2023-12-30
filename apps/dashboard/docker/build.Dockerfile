FROM mgarnier11/my-home:deps as dashboard-build

WORKDIR /build

COPY --from=mgarnier11/my-home:deps /build .
COPY ./apps/dashboard ./apps/dashboard
COPY --from=mgarnier11/my-home:deps /apps/dashboard/node_modules ./apps/dashboard/node_modules
COPY --from=mgarnier11/my-home:libs /build/libs ./libs

RUN pnpm run --filter dashboard build

