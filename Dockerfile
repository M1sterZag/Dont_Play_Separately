# syntax=docker/dockerfile:1

FROM golang:1.26.0-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/dps ./cmd/api

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY --from=build /bin/dps /usr/local/bin/dps

EXPOSE 5050

ENTRYPOINT ["dps"]
