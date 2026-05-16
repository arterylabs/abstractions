FROM golang:1.23-alpine AS build

RUN apk add --no-cache --update git build-base

WORKDIR /go/src/api

COPY . .

RUN go mod tidy && \
    go build -o /api cmd/main.go

FROM alpine

ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.7.3/wait /wait
RUN chmod +x /wait

COPY --from=build /api /api
COPY --from=build /go/src/api/db/migrations /db/migrations

RUN chmod +x /api

CMD /wait && /api
