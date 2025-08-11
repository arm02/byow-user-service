# Stage 1 - Build the app
FROM golang:1.24.5 as builder

WORKDIR /app

# Copy source
COPY . .

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g cmd/main.go

# Build the binary
RUN go build -o main cmd/main.go

# Stage 2 - Run the binary
FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/docs /app/docs

EXPOSE 8080

CMD ["/app/main"]
