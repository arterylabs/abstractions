FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /out/api ./cmd/main.go

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/api /api
COPY --from=builder /app/db/migrations /db/migrations

EXPOSE 1234

CMD ["/bin/sh", "-c", "/api migrate 0 && /api"]
