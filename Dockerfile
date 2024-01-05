# Use the official Golang image as the base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy only the necessary files into the container
COPY go.* ./
COPY cmd/ ./cmd

# Download dependencies
RUN go mod download

# Build the Go application
RUN go build -o main ./cmd

# Expose the port we will listen on
EXPOSE 8080

# Set environment variables for the database connection
ENV DB_USERNAME secret_username
ENV DB_PASSWORD secret_password
ENV DB_NAME secret_database_name

# Run the Go application
CMD ["./main"]