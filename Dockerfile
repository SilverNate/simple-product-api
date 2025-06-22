# Dockerfile
FROM golang:1.23-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go install github.com/joho/godotenv/cmd/godotenv@latest

COPY . .

RUN go build -o main ./cmd/main.go

EXPOSE 8080
CMD ["./main"]
