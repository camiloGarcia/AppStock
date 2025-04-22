# Etapa de build
FROM node:20-alpine AS builder
WORKDIR /app

# Recibir la variable de entorno
ARG VITE_API_BASE_URL
ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

COPY package*.json ./
RUN npm install

COPY . .

# 👇 Asegura que VITE_API_BASE_URL se use en build
RUN echo "VITE_API_BASE_URL=$VITE_API_BASE_URL" > .env && npm run build

# Etapa de producción
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
