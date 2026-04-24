# syntax=docker/dockerfile:1.7
# ---- Frontend build stage ------------------------------------------------
FROM node:20-alpine AS frontend
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci --no-audit --no-fund
COPY vite.config.js ./
COPY index.html ./
COPY public ./public
COPY src ./src
RUN npm run build

# ---- Backend build stage -------------------------------------------------
FROM golang:1.22-alpine AS backend
WORKDIR /src
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ ./
# Embed the freshly-built SPA into the Go binary.
COPY --from=frontend /src/dist ./dist
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o /out/che1-dashboard ./...

# ---- Runtime stage -------------------------------------------------------
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=backend /out/che1-dashboard /app/che1-dashboard
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/che1-dashboard"]
