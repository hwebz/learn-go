# Use the official Golang image
FROM golang:1.22-alpine

# Install dependencies
RUN apk add --no-cache git

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and go sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Install golang-migrate
RUN go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Copy migration files
COPY ./db/migration ./migrations

# Build the Go app
RUN go build -o go-simple-bank .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./go-simple-bank"]