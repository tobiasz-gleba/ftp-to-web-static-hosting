# Start with the official Go image
FROM golang:1.23.3

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# # Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY main.go main.go

# # # Build the Go application
RUN go build -o ftp-web-server.exe

# Expose the application's port
EXPOSE 80

# Run the executable
CMD ["./ftp-web-server.exe"]
