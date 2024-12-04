FROM golang:1.23-bookworm AS build

COPY ./ /go/src/github.com/softonic/hp-throttling

RUN cd /go/src/github.com/softonic/hp-throttling && make build

FROM scratch

COPY --from=build /go/src/github.com/softonic/hp-throttling/bin/hp-throttling /

ENTRYPOINT ["/hp-throttling", "-logtostderr"]
