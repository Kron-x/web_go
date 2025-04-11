# Compilating
FROM golang:1.22-alpine AS builder

WORKDIR /app_go

COPY app_go/ .

RUN ls /app_go -la  # Покажет файлы в /app_go

ENV CGO_ENABLED=0

RUN go mod download && \
    go build -o main ./cmd/web/

# running
FROM alpine:latest

WORKDIR /app_go

# Копируем скомпилированный бинарник из стадии builder
COPY --from=builder /app_go/main .

COPY app_go/configs/ ./configs/

COPY app_go/static/ ./static/

RUN mkdir -p ./logs

EXPOSE 5000 8080

CMD ["./main"]