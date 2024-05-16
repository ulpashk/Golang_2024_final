# Use the golang base image for the builder stage
FROM golang:1.22.1 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go module and sum files to the working directory
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the entire application directory into the container
COPY . .

# Build the application with CGO disabled for a Linux target
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/api

# Start a new stage from the alpine base image
FROM alpine:latest

# Install CA certificates to enable TLS in the Alpine image
RUN apk --no-cache add ca-certificates

# Set the working directory in the container
WORKDIR /root/

# Copy the compiled application from the builder stage
COPY --from=builder /app/app .

# Copy the migrations folder from the builder stage
COPY --from=builder /app/migrations ./migrations

# Define the command to run the executable
CMD ["./app"]