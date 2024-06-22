
FROM alpine:latest

WORKDIR /app

COPY listener-service /app
ENTRYPOINT ["/app/listener-service"]