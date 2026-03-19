# ── Stage 1: Build ────────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git (diperlukan go mod download untuk private modules)
RUN apk add --no-cache git

# Download dependencies dulu (cache layer — tidak rebuild jika kode berubah)
COPY go.mod go.sum ./
RUN go mod download

# Copy source dan build
COPY . .

ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /bin/bmt-api \
    ./cmd/server

# ── Stage 2: Runtime ──────────────────────────────────────────────────────────
# Distroless: tidak ada shell, tidak ada package manager — attack surface minimal
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /bin/bmt-api /bmt-api

EXPOSE 8080

ENTRYPOINT ["/bmt-api"]
