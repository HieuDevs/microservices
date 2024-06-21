
FROM alpine:latest

WORKDIR /app

COPY authentication-service /app
ENTRYPOINT ["/app/authentication-service"]