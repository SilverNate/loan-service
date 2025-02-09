# Build Stage: Use the latest Golang image to compile the application.
FROM golang:latest AS builder
WORKDIR /app

# Copy go.mod and go.sum, then download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source code and build the binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o loan_service cmd/app/main.go

# Final Stage: Use a minimal Alpine image to run the binary.
FROM alpine:latest
WORKDIR /root/

# Copy the statically compiled binary from the builder stage.
COPY --from=builder /app/loan_service .

# Expose the port the app listens on.
EXPOSE 8080

# Set the entrypoint to run the binary.
ENTRYPOINT ["./loan_service"]
