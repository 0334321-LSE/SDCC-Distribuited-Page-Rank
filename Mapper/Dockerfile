# syntax=docker/dockerfile:1
FROM golang:1.19.3

# Set destination for COPY
WORKDIR /Mapper

# Copy the source code
COPY . .

# Install net-tools to obtain open ports
RUN apt-get update && \
    apt-get install -y net-tools

# Download Go modules
RUN go mod download

# Build
RUN go build -o docker-mapper

# Run
CMD ["./docker-mapper"]