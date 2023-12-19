FROM golang:1.21-alpine AS builder

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download


# Copy the code into the container
COPY . .

# Build the application
# go build -o [name] [path to file]
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app cmd/main.go

# Move to /bin directory as the place for resulting binary folder
WORKDIR /bin

# Copy binary from build to main folder
RUN cp /build/app .


############################
# STEP 2 build a small image
############################
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY . .
COPY --from=builder /bin/app /
# Copy the code into the container

EXPOSE 5050

# Command to run the executable
ENTRYPOINT ["/app"]
