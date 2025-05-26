ARG GO_VERSION=1.24.0
FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache ca-certificates git gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

RUN go mod download && go mod verify

COPY . .

RUN go build \
    -a \
    -trimpath \
    -buildvcs=false \
    -ldflags="-s -w" \
    -installsuffix cgo \
    -o ./bin/go_api_server \
    ./cmd/api/main.go

FROM alpine:latest

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/bin/go_api_server /app/go_api_server

COPY --from=builder /app/internal/config/config.yaml /app/config.yaml

ENV TZ=Asia/Taipei

EXPOSE 8080

CMD ["/app/go_api_server"]