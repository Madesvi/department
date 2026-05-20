# --- Stage 1: Builder ---
FROM golang:1.25.7-alpine AS builder

RUN apk add --no-cache curl make nodejs npm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o main .

# --- Stage 2: Runner ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=builder /app/main .
ENV APP_VERSION=latest
EXPOSE 3000
CMD ["./main"]
