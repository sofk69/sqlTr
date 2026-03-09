FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o spend-management ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/spend-management .

COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080

CMD ["./spend-management"]