
FROM alpine:latest

WORKDIR /app

COPY broker-service /app
ENTRYPOINT ["/app/broker-service"]