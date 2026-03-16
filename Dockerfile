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

# Install Playwright dependencies and only Chromium browser
RUN apk add --no-cache \
    ca-certificates \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    nodejs \
    npm \
    curl \
    wget && \
  npm install -g playwright && \
  npx playwright install chromium && \
  npx playwright install-deps chromium && \
  npm cache clean --force && \
  rm -rf /root/.npm

# Set Playwright environment variables
ENV PLAYWRIGHT_BROWSERS_PATH=/root/.cache/ms-playwright
ENV PLAYWRIGHT_DRIVER_PATH=/usr/local/lib/node_modules/playwright

RUN go build -o certimate && \
  rm -rf /go/pkg


FROM alpine:latest

WORKDIR /app

# Install only runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    nodejs

# Set Playwright environment variables
ENV PLAYWRIGHT_BROWSERS_PATH=/root/.cache/ms-playwright
ENV PLAYWRIGHT_DRIVER_PATH=/usr/local/lib/node_modules/playwright

COPY --from=builder /app/certimate .
COPY --from=builder /root/.cache/ms-playwright /root/.cache/ms-playwright
COPY --from=builder /usr/local/lib/node_modules/playwright /usr/local/lib/node_modules/playwright

ENTRYPOINT ["./certimate", "serve", "--http", "0.0.0.0:8090"]