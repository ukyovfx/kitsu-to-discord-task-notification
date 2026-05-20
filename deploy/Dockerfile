FROM golang:1.21-bookworm AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o kitsu-discord ./src

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates curl tzdata && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/kitsu-discord .
CMD ["./kitsu-discord"]
