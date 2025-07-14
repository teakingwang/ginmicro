FROM golang:1.22-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o user-service ./cmd/user

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/user-service .
COPY resources/config.yaml /root/config.yaml

EXPOSE 50051
CMD ["./user-service"]
