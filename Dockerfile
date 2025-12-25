# Image de base Go légère
FROM golang:1.21-alpine

# Installer git (nécessaire pour les dépendances)
RUN apk add --no-cache git

# Définit le répertoire de travail
WORKDIR /app

# Copie les fichiers de dépendances d'abord (optimisation du cache Docker)
COPY go.mod go.sum ./

# Télécharge les dépendances
RUN go mod download

# Copie le code source
COPY . .

# Build l'application
RUN go build -o main ./cmd/api

# Expose le port
EXPOSE 8081

# Commande pour lancer l'application
CMD ["./main"]