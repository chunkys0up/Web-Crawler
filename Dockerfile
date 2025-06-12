# Use official Go image with the version you want
FROM golang:1.24

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to the container
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of your source code
COPY . .

# Build your Go application binary named "crawler"
RUN go build -o crawler

# Default command to run your app
CMD ["./crawler"]
