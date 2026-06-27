FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -ldflags="-s -w" -o doapi .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/doapi /doapi
EXPOSE 8080
ENTRYPOINT ["/doapi"]
