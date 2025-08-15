FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/app ./cmd/app

FROM alpine:3.18
RUN apk add --no-cache curl
COPY --from=builder /app/bin/app /app/app
COPY --from=builder /app/migrations /migrations
CMD ["/app/app"]