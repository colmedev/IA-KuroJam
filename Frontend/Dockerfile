# Build stage
FROM node:20-alpine3.17 AS builder

WORKDIR /app

# Copy package.json and pnpm-lock.yaml
COPY package.json pnpm-lock.yaml ./

# Install dependencies
RUN npm install -g pnpm && pnpm install

# Copy the rest of the application code
COPY . .

# Build the Vite application
ARG VITE_PUBLISHED_KEY
ARG VITE_BACKEND_URL
RUN VITE_PUBLISHED_KEY=$VITE_PUBLISHED_KEY \
    VITE_BACKEND_URL=$VITE_BACKEND_URL \
    pnpm build

# Serve stage
FROM nginx:stable-alpine

# Copy built assets to Nginx's web root
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy the custom Nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
