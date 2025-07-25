# use golang image
FROM golang:1.24-alpine3.21

# Install Git and Air
RUN apk add --no-cache git
RUN go install github.com/air-verse/air@latest

# set working dir
WORKDIR /app

# Copy go.mod and download deps (for cache)
COPY go.mod ./
COPY go.sum ./
RUN go mod download


# copy the source code
COPY . .


# export the port
EXPOSE 8000

# Run the executable
CMD ["air"]

