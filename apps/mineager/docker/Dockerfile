FROM build-go-app AS build-backend
FROM build-node-app AS build-frontend

FROM alpine AS runtime

WORKDIR /dist

COPY --from=build-backend /dist .
COPY --from=build-frontend /dist/app .

CMD ["/dist/go-proxy"]
