FROM golang:1.22

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o triple-s ./cmd/triple-s

EXPOSE 9000

CMD ["./triple-s", "--dir=/data"]
