FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o service-email main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/service-email .

EXPOSE 8084
CMD ["./service-email"]
