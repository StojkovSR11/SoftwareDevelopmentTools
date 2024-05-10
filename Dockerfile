# Build Stage
FROM golang:1.19 AS build

# Set working directory
WORKDIR /app

# Copy only Go module files
COPY go.mod go.sum ./

# Download Go modules (cached if go.mod/go.sum are unchanged)
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# Final Stage
FROM alpine:3.14

# Copy the binary from the build stage
COPY --from=build /main /main

# Expose port
EXPOSE 8000

# Run the binary
CMD ["/main"]
