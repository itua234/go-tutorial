# --- Stage 1: Builder ---
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker's cache.
# If these files don't change, the go mod download step will be skipped.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
# CGO_ENABLED=0 creates a statically linked binary (no external C dependencies)
# -o /app_name specifies the output path for the binary
RUN CGO_ENABLED=0 go build -o /app_name .


# --- Stage 2: Final / Production Image ---
# Use a minimal base image like 'alpine' or 'scratch' for the final image.
# Alpine is a good choice as it's very small and includes necessary tools like 'ca-certificates'.
FROM alpine:latest

# Install ca-certificates for secure HTTPS connections
# Use --no-cache to avoid storing package index, keeping the image small
RUN apk add --no-cache ca-certificates

# Set the working directory (optional, but good practice)
WORKDIR /root/

# Copy the compiled binary from the 'builder' stage
COPY --from=builder /app_name .

# Expose the port your application listens on (e.g., 8080)
EXPOSE 8080

# The command to run the application when the container starts
CMD ["./app_name"]