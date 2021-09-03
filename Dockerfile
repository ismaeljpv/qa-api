FROM golang:alpine

# Add Maintainer Info
LABEL maintainer="Ismael Pena <ismael.pena@bairesdev.com>"

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the code into the container
COPY . .

WORKDIR /app/cmd/questionary

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /app/cmd/questionary/main .

# Copy .env file
RUN cp /app/.env .

# Export necessary port
EXPOSE 8080

# Command to run when starting the container
CMD ["./main"]