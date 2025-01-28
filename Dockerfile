ARG GO_VERSION=1.22.4
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /mrt-app ./cmd/server

FROM debian:bookworm
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=builder /mrt-app /usr/local/bin/
CMD ["mrt-app"]