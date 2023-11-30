FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY ./go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .


RUN go build -o ./out ./cmd/main.go

# Path: Dockerfile
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/out .

EXPOSE 8080

CMD ["./out"]


