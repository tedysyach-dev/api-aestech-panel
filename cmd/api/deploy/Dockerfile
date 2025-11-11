# --- build stage ---
FROM golang:1.24-alpine AS build

WORKDIR /app

# Copy go.mod & go.sum dari root project
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy semua source dari root
COPY . .

# Build main.go di cmd/api
RUN go build -o /bin/app cmd/api/main.go

# --- runtime stage ---
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /root/

COPY --from=build /bin/app .
COPY cmd/api/deploy/config.json ./config.json

EXPOSE 8030

CMD ["./app"]
