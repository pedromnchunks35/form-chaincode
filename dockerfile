# Use an official Golang runtime as a parent image
FROM golang:1.24.0

# Set the working directory to /chaincode
WORKDIR /chaincode

# Copy the current directory contents into the container at /chaincode
COPY . /chaincode

RUN go mod download

# Build the chaincode
RUN go build -o run ./cmd/main.go

# Specify the command to run on container start
CMD ["./run"]
