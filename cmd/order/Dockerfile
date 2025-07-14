# cmd/order/Dockerfile
FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o order ./cmd/order

ENV SERVICE_ID=order-1
ENV SERVICE_NAME=order

EXPOSE 8080

CMD ["./order"]
