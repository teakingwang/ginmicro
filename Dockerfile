# Dockerfile
FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/${SERVICE}

EXPOSE 8080

ENTRYPOINT ["/app/app"]
