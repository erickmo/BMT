#!/bin/sh
# Ganti placeholder API_BASE_URL dengan nilai env var saat container start.
# Ini memungkinkan satu image Docker dipakai di berbagai environment
# (staging, production) tanpa perlu rebuild.

set -e

API_URL="${API_BASE_URL:-http://localhost:8080}"

echo "BMT Management — API_BASE_URL: ${API_URL}"

# Ganti semua placeholder di file JS hasil build Flutter
find /usr/share/nginx/html -name "*.js" -exec \
    sed -i "s|RUNTIME_API_BASE_URL|${API_URL}|g" {} +

# Jalankan nginx
exec nginx -g "daemon off;"
