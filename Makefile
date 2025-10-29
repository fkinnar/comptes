.PHONY: build clean test install help run sign

# Variables
BINARY_NAME=comptes
MAIN_PACKAGE=./cmd/comptes
BUILD_DIR=.
GO_FLAGS=-trimpath
LDFLAGS=-s -w

# Détection de l'OS
UNAME_S := $(shell uname -s)
IS_MACOS := $(shell uname -s | grep -q Darwin && echo "yes" || echo "no")

# Détection de Go (priorité à Homebrew sur macOS)
GO_CMD := go
ifeq ($(IS_MACOS),yes)
    ifeq ($(shell test -f /opt/homebrew/bin/go && echo "yes"),yes)
        GO_CMD := /opt/homebrew/bin/go
    endif
endif

help: ## Affiche l'aide
	@echo "Targets disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: ## Compile le binaire
	@echo "Compilation de $(BINARY_NAME)..."
	@if [ "$(IS_MACOS)" = "yes" ]; then \
		GO_VERSION=$$($(GO_CMD) version | awk '{print $$3}' | sed 's/go//'); \
		GO_MAJOR=$$(echo $$GO_VERSION | cut -d. -f1); \
		GO_MINOR=$$(echo $$GO_VERSION | cut -d. -f2); \
		if [ "$$GO_MAJOR" -lt 1 ] || ([ "$$GO_MAJOR" -eq 1 ] && [ "$$GO_MINOR" -lt 24 ]); then \
			echo "⚠️  ATTENTION: Go $$GO_VERSION détecté. macOS 26+ nécessite Go 1.24+"; \
			echo "   Lancez './update-go.sh' pour mettre à jour Go"; \
			echo "   OU consultez UPDATE_GO.md pour plus d'informations"; \
			echo ""; \
			echo "   Tentative de compilation (peut échouer avec 'missing LC_UUID'):"; \
		fi; \
	fi
	@$(GO_CMD) build $(GO_FLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@if [ "$(IS_MACOS)" = "yes" ]; then \
		echo "Signature du binaire pour macOS..."; \
		codesign --remove-signature $(BUILD_DIR)/$(BINARY_NAME) 2>/dev/null || true; \
		codesign --sign - $(BUILD_DIR)/$(BINARY_NAME) || true; \
	fi
	@echo "✓ Compilation réussie: $(BUILD_DIR)/$(BINARY_NAME)"

build-with-uuid: ## Compile avec LC_UUID (nécessite Go 1.24+ ou uuidgen)
	@echo "Compilation de $(BINARY_NAME) avec LC_UUID..."
	@$(GO_CMD) build $(GO_FLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@if [ "$(IS_MACOS)" = "yes" ]; then \
		if command -v uuidgen >/dev/null 2>&1; then \
			UUID=$$(uuidgen | tr '[:upper:]' '[:lower:]'); \
			echo "Ajout de LC_UUID (UUID: $$UUID)..."; \
			echo "Note: Cette méthode nécessite des outils supplémentaires. Mettez à jour Go vers 1.24+ pour une solution permanente."; \
		fi; \
		echo "Signature du binaire pour macOS..."; \
		codesign --remove-signature $(BUILD_DIR)/$(BINARY_NAME) 2>/dev/null || true; \
		codesign --sign - $(BUILD_DIR)/$(BINARY_NAME) || true; \
	fi
	@echo "✓ Compilation réussie: $(BUILD_DIR)/$(BINARY_NAME)"

build-debug: ## Compile avec informations de debug
	@echo "Compilation en mode debug..."
	@$(GO_CMD) build $(GO_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@if [ "$(IS_MACOS)" = "yes" ]; then \
		echo "Signature du binaire pour macOS..."; \
		codesign --remove-signature $(BUILD_DIR)/$(BINARY_NAME) 2>/dev/null || true; \
		codesign --sign - $(BUILD_DIR)/$(BINARY_NAME) || true; \
	fi
	@echo "✓ Compilation réussie: $(BUILD_DIR)/$(BINARY_NAME)"

clean: ## Nettoie les fichiers générés
	@echo "Nettoyage..."
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
	@$(GO_CMD) clean -cache -testcache
	@echo "✓ Nettoyage terminé"

test: ## Lance les tests
	@echo "Exécution des tests..."
	@$(GO_CMD) test -v ./...

test-coverage: ## Lance les tests avec couverture
	@echo "Exécution des tests avec couverture..."
	@$(GO_CMD) test -cover -coverprofile=coverage.out ./...
	@$(GO_CMD) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Rapport de couverture généré: coverage.html"

run: build ## Compile et exécute (utilise les arguments après --)
	@./$(BINARY_NAME) $(ARGS)

install: build ## Installe le binaire dans /usr/local/bin (nécessite sudo)
	@echo "Installation de $(BINARY_NAME) dans /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@echo "✓ Installation réussie"

fmt: ## Formate le code
	@echo "Formatage du code..."
	@$(GO_CMD) fmt ./...
	@echo "✓ Formatage terminé"

vet: ## Analyse statique du code
	@echo "Analyse statique..."
	@$(GO_CMD) vet ./...
	@echo "✓ Analyse terminée"

lint: fmt vet ## Formatage et analyse statique

mod-tidy: ## Nettoie les dépendances Go
	@echo "Nettoyage des dépendances..."
	@$(GO_CMD) mod tidy
	@echo "✓ Dépendances nettoyées"

mod-verify: ## Vérifie les dépendances Go
	@echo "Vérification des dépendances..."
	@$(GO_CMD) mod verify
	@echo "✓ Dépendances vérifiées"

sign: ## Signe le binaire (macOS uniquement)
	@if [ "$(IS_MACOS)" != "yes" ]; then \
		echo "Cette commande est uniquement disponible sur macOS"; \
		exit 1; \
	fi
	@if [ ! -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		echo "Erreur: $(BINARY_NAME) n'existe pas. Lancez 'make build' d'abord."; \
		exit 1; \
	fi
	@echo "Signature du binaire..."
	@codesign --remove-signature $(BUILD_DIR)/$(BINARY_NAME) 2>/dev/null || true
	@codesign --sign - $(BUILD_DIR)/$(BINARY_NAME)
	@echo "✓ Signature réussie"

# Cible par défaut
.DEFAULT_GOAL := build

