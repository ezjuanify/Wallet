FROM golang:1.24.4-alpine

WORKDIR /wallet

RUN apk add --no-cache git

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o wallet ./cmd/wallet

EXPOSE 8080

CMD ["./wallet"]