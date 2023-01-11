# Build stage
FROM golang:1.19.2-alpine3.16 AS builder
RUN apk update && apk add make
WORKDIR /app
COPY . .
RUN make build

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/bin/exporter .
COPY --from=builder /app/config config/

EXPOSE 8080
ENTRYPOINT [ "/app/exporter" ]