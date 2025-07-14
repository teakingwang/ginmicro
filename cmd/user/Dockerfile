# cmd/user/Dockerfile
FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o user ./cmd/user

ENV SERVICE_ID=user-1
ENV SERVICE_NAME=user


EXPOSE 8080

CMD ["./user"]
