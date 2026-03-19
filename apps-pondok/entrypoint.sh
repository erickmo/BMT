#!/bin/sh
set -e

API_URL="${API_BASE_URL:-http://localhost:8080}"

echo "BMT Pondok — API_BASE_URL: ${API_URL}"

find /usr/share/nginx/html -name "*.js" -exec \
    sed -i "s|RUNTIME_API_BASE_URL|${API_URL}|g" {} +

exec nginx -g "daemon off;"
