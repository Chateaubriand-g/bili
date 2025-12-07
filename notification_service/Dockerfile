FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY . .
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go get github.com/Chateaubriand-g/bili/common && \
    go mod tidy

RUN CGO_ENABLED=0 go build -ldflags "-w -s" -o notify_service main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/notify_service ./
COPY --from=builder /app/config ./config

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

RUN chmod +x ./notify_service && \
    chown -R appuser:appgroup /app

USER appuser
EXPOSE 8083

CMD ["./notify_service"]
