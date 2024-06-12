# Stage 1: Build the Go binaries
FROM golang:1.22 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the workspace
COPY go.mod ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go client binary
RUN CGO_ENABLED=0 go build -o client ./cmd/eftepcli

# Build the Go eftepd binary
RUN CGO_ENABLED=0 go build -o eftepd ./cmd/eftepd

# Stage 2: Run the Go binaries
FROM debian:stable-slim

# Set the working directory inside the container
WORKDIR /app

# Copy the binaries from the builder stage
COPY --from=builder /app/client /app/eftepcli
COPY --from=builder /app/eftepd /app/eftepd

# Expose the port the app runs on, if necessary
EXPOSE 8080 8081

# Set the entrypoint for the container
ENTRYPOINT ["sh", "-c"]
CMD ["./eftepd"]
