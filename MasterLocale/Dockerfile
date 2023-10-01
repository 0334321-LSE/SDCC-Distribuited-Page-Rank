# syntax=docker/dockerfile:1
FROM golang:1.19.3

# Set destination for COPY
WORKDIR /Master

# Copy the source code
COPY . .

# Download Go modules
RUN go mod download

# Build
RUN go build -o docker-master

# Run
CMD ["./docker-master"]