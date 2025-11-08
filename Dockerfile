# multi-stage build: compile with official golang image, run in small alpine image
FROM golang:1.25.3-alpine AS builder

WORKDIR /src

# cache modules
COPY go.mod go.sum ./
RUN apk add --no-cache git
RUN go mod download

# copy, build
COPY . .
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -ldflags="-s -w" -o /app/go-finance .

# runtime image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /app/go-finance .

EXPOSE 8080

ENTRYPOINT ["/app/go-finance"]