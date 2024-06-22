
FROM alpine:latest

WORKDIR /app

COPY logger-service /app
ENTRYPOINT ["/app/logger-service"]