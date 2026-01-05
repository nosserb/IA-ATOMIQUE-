.PHONY: help build run web test clean docker-build docker-run

# Variables
BINARY_NAME=programme
PORT=8080
GO=go
DOCKER_IMAGE=ia-atomique
DOCKER_TAG=latest

help:
	@echo "╔════════════════════════════════════════╗"
	@echo "║  IA-ATOMIQUE Makefile                  ║"
	@echo "╚════════════════════════════════════════╝"
	@echo ""
	@echo "Commandes disponibles:"
	@echo ""
	@echo "  make build          - Compiler le projet"
	@echo "  make web            - Lancer l'interface web"
	@echo "  make run            - Lancer en mode interactif"
	@echo "  make test           - Tester la compilation"
	@echo "  make clean          - Supprimer les fichiers compilés"
	@echo "  make docker-build   - Construire l'image Docker"
	@echo "  make docker-run     - Lancer le conteneur Docker"
	@echo "  make deps           - Mettre à jour les dépendances"
	@echo ""

build:
	@echo "📦 Compilation du projet..."
	$(GO) build -o $(BINARY_NAME)
	@echo "✓ Compilation terminée: ./$(BINARY_NAME)"

web: build
	@echo "🌐 Démarrage de l'interface web..."
	@echo "   URL: http://localhost:$(PORT)"
	./$(BINARY_NAME) web

run: build
	@echo "💬 Mode interactif..."
	./$(BINARY_NAME) interactive

test:
	@echo "🧪 Test de compilation..."
	$(GO) build -o $(BINARY_NAME)
	@echo "✓ Code compilé avec succès"
	@echo "🌐 Test du serveur (5 secondes)..."
	@timeout 5 ./$(BINARY_NAME) web 2>&1 | head -20 || true
	@echo "✓ Tests terminés"

clean:
	@echo "🧹 Nettoyage..."
	$(GO) clean
	rm -f $(BINARY_NAME)
	@echo "✓ Nettoyage terminé"

deps:
	@echo "📚 Mise à jour des dépendances..."
	$(GO) mod tidy
	$(GO) mod download
	@echo "✓ Dépendances mises à jour"

docker-build:
	@echo "🐳 Construction de l'image Docker..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "✓ Image construite: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-run: docker-build
	@echo "🚀 Lancement du conteneur..."
	@echo "   URL: http://localhost:$(PORT)"
	docker run -p $(PORT):8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-shell: docker-build
	docker run -it -p $(PORT):8080 $(DOCKER_IMAGE):$(DOCKER_TAG) /bin/sh

all: clean deps build test
	@echo "✅ Build complet réussi!"
