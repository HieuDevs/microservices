
FROM alpine:latest

WORKDIR /app

COPY mail-service /app
COPY templates /app/templates
ENTRYPOINT ["/app/mail-service"]