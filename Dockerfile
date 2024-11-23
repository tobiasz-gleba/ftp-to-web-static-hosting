# Step 1: Start with the official Go image
FROM golang:1.23.3 as build

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY main.go main.go

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o ftp-web-server.exe

# Step 2: Runtime stage
FROM scratch

# Copy only the binary from the build stage to the final image
COPY --from=build /app/ftp-web-server.exe /

# Expose the application's port
EXPOSE 80

# Run the executable
ENTRYPOINT ["/ftp-web-server.exe"]
