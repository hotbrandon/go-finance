# multi-stage build: compile with official golang image, run in small alpine image
FROM golang:1.21-alpine AS builder

WORKDIR /src

# bring in modules first for caching
COPY go.mod go.sum ./
RUN apk add --no-cache git
RUN go mod download

# copy the rest of the sources and build
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o /app/go-finance ./


# runtime image
FROM alpine:3.18
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

# copy binary from builder
COPY --from=builder /app/go-finance .

# default port and env (can be overridden in docker-compose or at runtime)
EXPOSE 8080

# run the binary
ENTRYPOINT ["/app/go-finance"]