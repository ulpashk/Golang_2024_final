# Step 1: Use the official Golang image as a builder.
FROM golang:1.16 as builder

# Step 2: Set the Current Working Directory inside the container
WORKDIR /app

# Step 3: Copy the Go Module files
COPY go.mod go.sum ./

# Step 4: Download all dependencies.
RUN go mod download

# Step 5: Copy the source code into the container
COPY . .

# Step 6: Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd/api

# Step 7: Start from a smaller image
FROM alpine:latest  

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/myapp .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]

