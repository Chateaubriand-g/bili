FROM go:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags "-w -s" -o gateway main.go


FROM alpine:latest

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/gateway ./
COPY --from=builder /app/config /app/config

RUN chmod +x ./gateway && \
    chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

CMD ["./gateway"]
