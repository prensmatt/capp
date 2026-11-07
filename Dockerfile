FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency definition files first for optimal caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the backend source code
COPY . .

# Build directly pointing to the main.go file inside cmd
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Final lightweight runner container
FROM alpine:latest
WORKDIR /app
RUN mkdir -p static/images static/songs
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static

EXPOSE 8080
CMD ["/app/main"]