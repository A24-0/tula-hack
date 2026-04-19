FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o server ./cmd/server

FROM alpine:3.19
RUN apk add --no-cache ca-certificates wget
WORKDIR /app
COPY --from=builder /app/server .
CMD ["./server"]
