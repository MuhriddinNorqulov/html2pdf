FROM golang:1.25.1 AS builder

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY ./internal ./internal
COPY ./cmd ./cmd

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/main.go


FROM debian:12-slim

ENV CONTAINER_MODE=1

RUN apt-get update \
 && apt-get install -y --no-install-recommends \
    wkhtmltopdf \
    fontconfig \
    fonts-dejavu \
    fonts-liberation \
    fonts-noto \
    fonts-noto-cjk \
    fonts-noto-color-emoji \
    ca-certificates \
    tzdata \
 && rm -rf /var/lib/apt/lists/* \
 && fc-cache -f


WORKDIR /app
COPY --from=builder /out/app ./app

ENTRYPOINT ["./app"]