# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS builder

WORKDIR /build

# Installer les dépendances nécessaires
RUN apk add --no-cache git

# Copier les fichiers du projet
COPY . .

# Compiler l'application
RUN go build -o programme .

# Stage final
FROM alpine:latest

WORKDIR /app

# Installer les dépendances runtime si nécessaire
RUN apk add --no-cache ca-certificates

# Copier le binaire compilé
COPY --from=builder /build/programme /app/programme

# Copier les fichiers web
COPY --from=builder /build/web /app/web

# Exposer le port du serveur web
EXPOSE 8080

# Lancer l'application
CMD ["./programme", "web"]
