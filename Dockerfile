FROM golang:1.23-alpine AS build

ENV GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd/
COPY internal internal/

RUN go build -o /app/rill ./cmd/rill/main.go

FROM alpine:3.20

WORKDIR /app

RUN adduser -D rill && chown -R rill:rill /app
USER rill

COPY --chown=rill:rill --chmod=440 .env.* /app/
COPY --from=build --chown=rill:rill --chmod=770 /app/rill /app/

EXPOSE 8910

ENTRYPOINT ["./rill"]
