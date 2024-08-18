FROM build AS build


FROM scratch AS runtime

WORKDIR /dist

COPY --from=build /dist .

RUN ls -la; \
  ls -la /dist

CMD ["/dist/go-proxy"]