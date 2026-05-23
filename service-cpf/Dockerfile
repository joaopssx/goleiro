FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o service-cpf main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/service-cpf .

EXPOSE 8081
CMD ["./service-cpf"]
