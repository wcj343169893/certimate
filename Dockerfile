FROM node:20-alpine AS front-builder

WORKDIR /app

COPY . /app/

RUN \
  cd /app/ui && \
  npm install && \
  npm run build


FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY ../. /app/

RUN rm -rf /app/ui/dist

COPY --from=front-builder /app/ui/dist /app/ui/dist

RUN go build -o certimate



FROM debian:bookworm-slim

WORKDIR /app

# 安装 Playwright 运行所需的系统依赖
RUN apt-get update && apt-get install -y \
    ca-certificates \
    libnss3 \
    libnspr4 \
    libatk1.0-0 \
    libatk-bridge2.0-0 \
    libcups2 \
    libdrm2 \
    libxkbcommon0 \
    libxcomposite1 \
    libxdamage1 \
    libxext6 \
    libxfixes3 \
    libxrandr2 \
    libgbm1 \
    libpango-1.0-0 \
    libcairo2 \
    libasound2 \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/certimate .

ENTRYPOINT ["./certimate", "serve", "--http", "0.0.0.0:8090"]
