FROM golang:1.24-alpine as build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o ./bin/api ./main.go

EXPOSE 3000

CMD ["./bin/api"]

