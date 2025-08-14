# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /bin/shazam ./main.go

# Final image
FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /bin/shazam /bin/shazam
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/bin/shazam"]
