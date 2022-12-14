# Start from the latest golang base image
FROM golang:alpine AS builder

RUN apk add git

ENV GO111MODULE=on

# Set the Current Working Directory inside the container
WORKDIR /app

RUN go env -w GOPROXY=https://goproxy.cn,https://goproxy.io,direct

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go env && go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
WORKDIR /app/cmd/concurrent-http-server
RUN go build -o /concurrent-http-server

FROM alpine:latest

WORKDIR /app/

COPY --from=builder /concurrent-http-server .

RUN mkdir /var/files/

EXPOSE 8080

ENTRYPOINT ["./concurrent-http-server"]

# Run server
CMD ["run"]
