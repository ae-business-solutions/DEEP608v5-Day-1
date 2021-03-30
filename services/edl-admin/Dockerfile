FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"'

FROM scratch
COPY --from=builder /go/bin/edl-admin /edl-admin
COPY --from=builder /go/src/app/assets/ /assets/
ENTRYPOINT ["/edl-admin"]