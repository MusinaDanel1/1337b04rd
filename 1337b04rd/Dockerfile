FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o 1337b04rd ./cmd/1337b04rd

EXPOSE 8080

CMD ["./1337b04rd"]
