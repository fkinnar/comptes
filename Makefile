.PHONY: build clean test install help run sign

# Variables
BINARY_NAME=comptes
MAIN_PACKAGE=./cmd/comptes
BUILD_DIR=.
GO_FLAGS=-trimpath
LDFLAGS=-s -w

# Détection de l'OS (compatible Windows/Unix)
# Vérifie si on est sur Windows en testant COMSPEC ou OS
ifdef OS
    ifeq ($(OS),Windows_NT)
        IS_WINDOWS := yes
        IS_MACOS := no
    else
        IS_WINDOWS := no
        IS_MACOS := $(shell uname -s 2>/dev/null | grep -q Darwin && echo "yes" || echo "no")
    endif
else ifdef COMSPEC
    IS_WINDOWS := yes
    IS_MACOS := no
else
    # Unix/Linux/macOS - teste uname
    UNAME_S := $(shell uname -s 2>/dev/null)
    ifeq ($(UNAME_S),Darwin)
        IS_WINDOWS := no
        IS_MACOS := yes
    else
        IS_WINDOWS := no
        IS_MACOS := no
    endif
endif

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
ifeq ($(IS_MACOS),yes)
	@GO_VERSION=$$($(GO_CMD) version | awk '{print $$3}' | sed 's/go//'); \
	GO_MAJOR=$$(echo $$GO_VERSION | cut -d. -f1); \
	GO_MINOR=$$(echo $$GO_VERSION | cut -d. -f2); \
	if [ "$$GO_MAJOR" -lt 1 ] || ([ "$$GO_MAJOR" -eq 1 ] && [ "$$GO_MINOR" -lt 24 ]); then \
		echo "⚠️  ATTENTION: Go $$GO_VERSION détecté. macOS 26+ nécessite Go 1.24+"; \
		echo "   Lancez './update-go.sh' pour mettre à jour Go"; \
		echo "   OU consultez UPDATE_GO.md pour plus d'informations"; \
		echo ""; \
		echo "   Tentative de compilation (peut échouer avec 'missing LC_UUID'):"; \
	fi
endif
	@$(GO_CMD) build $(GO_FLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)$(if $(filter yes,$(IS_WINDOWS)),.exe,) $(MAIN_PACKAGE)
ifeq ($(IS_MACOS),yes)
	@echo "Signature du binaire pour macOS..."; \
	codesign --remove-signature $(BUILD_DIR)/$(BINARY_NAME) 2>/dev/null || true; \
	codesign --sign - $(BUILD_DIR)/$(BINARY_NAME) || true
endif
	@echo "✓ Compilation réussie: $(BUILD_DIR)/$(BINARY_NAME)$(if $(filter yes,$(IS_WINDOWS)),.exe,)"

build-with-uuid: ## Compile avec LC_UUID (nécessite Go 1.24+ ou uuidgen)
	@echo "Compilation de $(BINARY_NAME) avec LC_UUID..."
	@$(GO_CMD) build $(GO_FLAGS) -ldflags="$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)$(if $(filter yes,$(IS_WINDOWS)),.exe,) $(MAIN_PACKAGE)
ifeq ($(IS_MACOS),yes)
	@if command -v uuidgen >/dev/null 2>&1; then \
		UUID=$$(uuidgen | tr '[:upper:]' '[:lower:]'); \
		echo "Ajout de LC_UUID (UUID: $$UUID)..."; \
		echo "Note: Cette méthode nécessite des outils supplémentaires. Mettez à jour Go vers 1.24+ pour une solution permanente."; \
	fi; \
	echo "Signature du binaire pour macOS..."; \
	codesign --remove-signature $(BUILD_DIR)/$(BINARY_NAME) 2>/dev/null || true; \
	codesign --sign - $(BUILD_DIR)/$(BINARY_NAME) || true
endif
	@echo "✓ Compilation réussie: $(BUILD_DIR)/$(BINARY_NAME)$(if $(filter yes,$(IS_WINDOWS)),.exe,)"

build-debug: ## Compile avec informations de debug
	@echo "Compilation en mode debug..."
	@$(GO_CMD) build $(GO_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)$(if $(filter yes,$(IS_WINDOWS)),.exe,) $(MAIN_PACKAGE)
ifeq ($(IS_MACOS),yes)
	@echo "Signature du binaire pour macOS..."; \
	codesign --remove-signature $(BUILD_DIR)/$(BINARY_NAME) 2>/dev/null || true; \
	codesign --sign - $(BUILD_DIR)/$(BINARY_NAME) || true
endif
	@echo "✓ Compilation réussie: $(BUILD_DIR)/$(BINARY_NAME)$(if $(filter yes,$(IS_WINDOWS)),.exe,)"

clean: ## Nettoie les fichiers générés
	@echo "Nettoyage..."
ifeq ($(IS_WINDOWS),yes)
	@if exist $(BUILD_DIR)\$(BINARY_NAME).exe del /q $(BUILD_DIR)\$(BINARY_NAME).exe
	@if exist $(BUILD_DIR)\$(BINARY_NAME) del /q $(BUILD_DIR)\$(BINARY_NAME)
else
	@rm -f $(BUILD_DIR)/$(BINARY_NAME)
endif
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
ifeq ($(IS_WINDOWS),yes)
	@$(BUILD_DIR)\$(BINARY_NAME).exe $(ARGS)
else
	@./$(BINARY_NAME) $(ARGS)
endif

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
ifeq ($(IS_MACOS),yes)
	@if [ ! -f $(BUILD_DIR)/$(BINARY_NAME) ]; then \
		echo "Erreur: $(BINARY_NAME) n'existe pas. Lancez 'make build' d'abord."; \
		exit 1; \
	fi
	@echo "Signature du binaire..."
	@codesign --remove-signature $(BUILD_DIR)/$(BINARY_NAME) 2>/dev/null || true
	@codesign --sign - $(BUILD_DIR)/$(BINARY_NAME)
	@echo "✓ Signature réussie"
else
	@echo "Cette commande est uniquement disponible sur macOS"
	@exit 1
endif

# Cible par défaut
.DEFAULT_GOAL := build
