FROM golang:1.15.5-alpine as build

WORKDIR /go/src/just-tit
ADD go.mod .
RUN apk add git gcc libc-dev ca-certificates

# Recompile the standard library with CGO
RUN CGO_ENABLED=1 go install -a std

ADD . .
# Compile the binary and statically link
RUN CGO_ENABLED=1 go build -ldflags '-linkmode external -extldflags -static' -o just-tit

FROM scratch
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/just-tit/just-tit /just-tit
COPY --from=build /go/src/just-tit/static/ /static/
COPY --from=build /go/src/just-tit/views/ /views/
COPY --from=build /go/src/just-tit/conf/ /conf/

CMD ["/just-tit"]
EXPOSE 8080
HEALTHCHECK --interval=5m --timeout=3s CMD /bin/true