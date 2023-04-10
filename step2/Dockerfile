# Use the official Go image as the base image
FROM golang:1.16-alpine

# Set the working directory
WORKDIR /app

# Copy the Go service code into the container
COPY main.go .
COPY go.mod .
COPY go.sum .
COPY index.html .
RUN mkdir ./gifs



# Build the Go service
RUN go build -o app

# Expose the port the service listens on
EXPOSE 8889

# Start the service
CMD ["./app"]
