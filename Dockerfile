# Stage 1: Build frontend
FROM node:22-alpine AS frontend-build
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# Stage 2: Build Go binary
FROM golang:1.24-alpine AS go-build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-build /app/frontend/dist ./internal/web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o /iplayer-arr ./cmd/iplayer-arr/

# Stage 3: Runtime
FROM alpine:3.21
RUN apk add --no-cache ffmpeg tzdata
COPY --from=go-build /iplayer-arr /usr/local/bin/iplayer-arr

EXPOSE 8191
VOLUME ["/config", "/downloads"]

ENV TZ=Europe/London

ENTRYPOINT ["iplayer-arr"]
