# Use the official Go image as the base image
FROM golang:1.16-alpine

# Set the working directory
WORKDIR /app

# Copy the Go service code into the container
COPY main.go .
COPY index.html .


# Build the Go service
RUN go build -o app

# Expose the port the service listens on
EXPOSE 8888

# Start the service
CMD ["./app"]
