FROM golang:1.24-alpine AS build

WORKDIR /api

# Install required packages
RUN apk add --no-cache git

# Installation des certificats CA, cruciaux pour les appels HTTPS.
RUN apk add --no-cache ca-certificates

# Installation de Swag CLI pour la documentation Swagger.
#RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod and sum files first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Génération de la documentation Swagger.
#RUN swag init --generalInfo ./routes/swagger_routes.go --output ./docs

# Configuration des variables d'environnement nécessaires pour la construction.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Construction du serveur API.
RUN go build -v -o apiserver .

# Utilisation d'une image scratch pour un environnement d'exécution minimal.
FROM scratch

# Copie du binaire, des fichiers de configuration, et de la documentation Swagger.
COPY --from=build /api/apiserver /
COPY --from=build /api/.env /
#COPY --from=build /api/docs /docs

# Copie des certificats CA.
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Command for verifying the binary
CMD ["./apiserver"]
