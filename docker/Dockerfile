FROM golang:1.18-alpine AS builder

RUN apk update && \
    apk add make git bash curl openssl alpine-sdk --no-cache

WORKDIR /app

COPY . .

RUN make build

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/main /app/main

RUN chmod +x main

CMD ["./main"]