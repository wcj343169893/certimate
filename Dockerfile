FROM node:20-alpine3.19 AS front-builder

WORKDIR /app

COPY . /app/

RUN \
  cd /app/ui && \
  npm install --no-audit --no-fund && \
  npm run build && \
  rm -rf node_modules


FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . /app/

RUN rm -rf /app/ui/dist

COPY --from=front-builder /app/ui/dist /app/ui/dist

# Install base dependencies
RUN apk update && \
  apk add --no-cache \
    ca-certificates \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    nodejs \
    npm

# Install curl and wget separately
RUN apk add --no-cache \
    curl \
    wget

# Install Playwright
RUN npm install -g playwright

# Install only Chromium browser
RUN npx playwright install chromium

# Install Chromium dependencies manually instead of using playwright install-deps
RUN apk add --no-cache \
    libstdc++ \
    glib \
    libx11 \
    libxcomposite \
    libxdamage \
    libxext \
    libxfixes \
    libxrandr \
    libxrender \
    libxtst \
    libxcb \
    libxkbcommon \
    mesa-gl \
    dbus \
    tzdata

# Clean up
RUN npm cache clean --force && \
  rm -rf /root/.npm && \
  rm -rf /var/cache/apk/*

# Set Playwright environment variables
ENV PLAYWRIGHT_BROWSERS_PATH=/root/.cache/ms-playwright
ENV PLAYWRIGHT_DRIVER_PATH=/usr/local/lib/node_modules/playwright

# Build the application
RUN go build -o certimate && \
  rm -rf /go/pkg


FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk update && \
  apk add --no-cache \
    ca-certificates \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    nodejs && \
  rm -rf /var/cache/apk/*

# Set Playwright environment variables
ENV PLAYWRIGHT_BROWSERS_PATH=/root/.cache/ms-playwright
ENV PLAYWRIGHT_DRIVER_PATH=/usr/local/lib/node_modules/playwright

COPY --from=builder /app/certimate .
COPY --from=builder /root/.cache/ms-playwright /root/.cache/ms-playwright
COPY --from=builder /usr/local/lib/node_modules/playwright /usr/local/lib/node_modules/playwright

ENTRYPOINT ["./certimate", "serve", "--http", "0.0.0.0:8090"]