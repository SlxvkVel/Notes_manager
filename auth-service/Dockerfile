FROM golang:1.24.3-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -ldflags="-w -s" -o main ./cmd/server

EXPOSE 8080

CMD ["./main"]