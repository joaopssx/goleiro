FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o service-cnpj main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/service-cnpj .

EXPOSE 8082
CMD ["./service-cnpj"]
