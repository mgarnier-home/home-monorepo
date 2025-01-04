FROM build-go-app AS build-backend
FROM build-node-app AS build-frontend

FROM alpine AS runtime

WORKDIR /dist

COPY --from=build-backend /dist .
COPY --from=build-frontend /dist/front .

ENV FRONTEND_PATH=/dist/browser
ENV SERVER_PORT=8080

CMD ["/dist/mineager-back"]
