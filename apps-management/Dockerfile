# ── Stage 1: Build Flutter Web ────────────────────────────────────────────────
FROM ghcr.io/cirruslabs/flutter:stable AS builder

WORKDIR /app

# Download dependencies dulu (cache layer)
COPY pubspec.yaml pubspec.lock ./
RUN flutter pub get

# Copy source
COPY . .

# Build dengan placeholder — API_BASE_URL diganti di runtime via entrypoint.sh
RUN flutter build web --release \
    --dart-define=API_BASE_URL=RUNTIME_API_BASE_URL \
    --dart-define=APP_NAME=BMT-Management

# ── Stage 2: Nginx ────────────────────────────────────────────────────────────
FROM nginx:1.27-alpine

COPY --from=builder /app/build/web /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

EXPOSE 80

ENTRYPOINT ["/entrypoint.sh"]
