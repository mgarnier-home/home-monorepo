FROM build-go-app AS build

FROM alpine:latest AS runtime

# Install curl for health check
RUN apk --no-cache add curl

WORKDIR /dist

COPY --from=build /dist .

CMD ["/dist/go-autosaver"]

HEALTHCHECK --interval=30s --timeout=10s --retries=3 \
  CMD wget -q --spider --timeout=5 --tries=1 http://localhost:8080 || exit 1
