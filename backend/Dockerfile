FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Generate go.sum and build
RUN go mod tidy && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Build the seed binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o seed ./cmd/seed

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binaries from builder
COPY --from=builder /app/main .
COPY --from=builder /app/seed .

EXPOSE 8080

CMD ["./main"]

