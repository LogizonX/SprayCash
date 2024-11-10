FROM golang:1.22-alpine AS build

run apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod download

# build the app

RUN go build -o api cmd/main.go

# stage 2 build an image to run the app
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=build /app/api .

# Expose the application's port (change as needed)
EXPOSE 8080

# Run the application
CMD ["./api"]