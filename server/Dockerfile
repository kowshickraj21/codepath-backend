FROM golang:1.23.2-alpine as builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main main.go

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /app/main .
COPY --from=builder /app/problems .

CMD ["./main"]
