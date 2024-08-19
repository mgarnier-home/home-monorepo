FROM build AS build


FROM scratch AS runtime

WORKDIR /dist

COPY --from=build /dist .

CMD ["/dist/go-proxy"]