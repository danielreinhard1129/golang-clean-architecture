# ---- Build Stage ----
# Use the official Go image that matches the version in go.mod (go 1.25.1)
# Using the -alpine variant to keep the build stage image small.
FROM golang:1.25-alpine AS builder

# Install root certificates and git. Git is required for 'go mod download'
# if there are dependencies directly from git repositories.
RUN apk add --no-cache ca-certificates git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies first
# This leverages Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application.
# CGO_ENABLED=0 is important for creating a statically linked binary 
# compatible with alpine/scratch images.
# -ldflags="-s -w" removes debug symbols to reduce binary size.
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags="-w -s" -o /app/server ./cmd/app/main.go

# ---- Final Stage ----
# Use a minimal base image that already includes root certificates.
# 'alpine' is an excellent choice, only ~5MB.
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the build stage
COPY --from=builder /app/server .

# Copy the .env file and configs directory required by the app at runtime
COPY .env .

# Expose the port your application uses (e.g., 8000)
EXPOSE 8000

# Command to run your application
CMD ["./server"]
