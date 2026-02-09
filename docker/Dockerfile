FROM golang:1.25.1 AS builder

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY ./internal ./internal
COPY ./cmd ./cmd

RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/main.go


FROM surnet/alpine-wkhtmltopdf:3.23.2-024b2b2-full

ENV CONTAINER_MODE=1

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /out/app ./app

ENTRYPOINT ["./app"]