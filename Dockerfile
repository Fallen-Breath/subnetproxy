FROM golang:1.24 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 go build -o subnetproxy cmd/subnetproxy/subnetproxy.go

FROM alpine:latest
RUN apk add --no-cache tzdata  # so TZ environment works

COPY --from=builder /build/subnetproxy /usr/bin/subnetproxy

ENTRYPOINT ["/usr/bin/subnetproxy"]
