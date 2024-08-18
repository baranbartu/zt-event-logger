# Stage 1: Build
FROM golang:1.21.3 AS builder

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o app

# Stage 2: Run
FROM golang:1.21.3 AS runner

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/app .

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./app"]
