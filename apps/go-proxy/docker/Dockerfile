FROM build-go-app AS build


FROM alpine:latest AS runtime

# Install curl for health check
RUN apk --no-cache add curl

WORKDIR /dist

COPY --from=build /dist .

CMD ["/dist/go-proxy"]

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:3000 || exit 1
