FROM node:24.22.2-alpine AS builder
WORKDIR /app
COPY web/shortener/package.json web/shortener/package-lock.json ./
RUN npm ci
COPY web/shortener/ ./
RUN npm run build

FROM nginx:1.30.6-alpine
COPY configs/nginx/web.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
