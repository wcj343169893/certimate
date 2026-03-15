FROM node:20-alpine3.19 AS front-builder

WORKDIR /app

COPY . /app/

RUN \
  cd /app/ui && \
  npm install && \
  npm run build


FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY . /app/

RUN rm -rf /app/ui/dist

COPY --from=front-builder /app/ui/dist /app/ui/dist

# Install Playwright dependencies
RUN apk add --no-cache \
    ca-certificates \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    nodejs \
    npm

# Install Playwright
RUN npm install -g playwright
RUN npx playwright install

RUN go build -o certimate


FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ttf-freefont \
    nodejs

COPY --from=builder /app/certimate .
COPY --from=builder /root/.cache/ms-playwright-go /root/.cache/ms-playwright-go
COPY --from=builder /root/.cache/ms-playwright /root/.cache/ms-playwright

ENTRYPOINT ["./certimate", "serve", "--http", "0.0.0.0:8090"]