FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project into the container
COPY . .

# Download dependencies
RUN go mod download

# Use `go run` to execute the project
CMD ["go", "run", "."]