FROM golang:1.18-alpine

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

EXPOSE 8081

CMD ["/app/main"]