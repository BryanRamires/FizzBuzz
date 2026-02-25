# syntax=docker/dockerfile:1

############################
# ---- Build stage ----
############################
FROM golang:1.25.6-alpine3.23 AS build

# Needed for HTTPS when downloading modules
RUN apk add --no-cache ca-certificates

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build static binary (required for distroless)
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -ldflags="-s -w" \
    -o /out/fizzbuzz-api ./cmd/api


############################
# ---- Runtime stage ----
############################
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy only the binary
COPY --from=build /out/fizzbuzz-api /fizzbuzz-api

# Document listening port
EXPOSE 8090

HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
  CMD ["/fizzbuzz-api", "-healthcheck"]

# Run as non-root (already default in this image)
ENTRYPOINT ["/fizzbuzz-api"]