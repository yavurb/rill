FROM golang:1.23-alpine AS build

ENV GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd/
COPY internal internal/
COPY config config/

RUN go build -o /app/rill ./cmd/rill/main.go

FROM alpine:3.20

WORKDIR /app

ARG PKL_VERSION=0.26.3

RUN apk add --no-cache curl \
  && curl -L -o /usr/local/bin/pkl "https://github.com/apple/pkl/releases/download/${PKL_VERSION}/pkl-alpine-linux-amd64" \
  && chmod +x /usr/local/bin/pkl

RUN adduser -D rill && chown -R rill:rill /app
USER rill

ARG ENVIRONMENT=production

COPY --chown=rill:rill --chmod=440 config/Config.pkl /app/config/
COPY --chown=rill:rill --chmod=440 config/${ENVIRONMENT}/config.pkl /app/config/${ENVIRONMENT}/
COPY --from=build --chown=rill:rill --chmod=770 /app/rill /app/

EXPOSE 8910

ENTRYPOINT ["./rill"]
