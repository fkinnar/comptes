#!/bin/bash
# Script pour mettre à jour Go vers une version compatible avec macOS 26+

set -e

echo "🔍 Vérification de la version actuelle de Go..."
CURRENT_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "   Version actuelle: $CURRENT_VERSION"

GO_MAJOR=$(echo $CURRENT_VERSION | cut -d. -f1)
GO_MINOR=$(echo $CURRENT_VERSION | cut -d. -f2)

if [ "$GO_MAJOR" -gt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -ge 24 ]); then
    echo "✅ Votre version de Go ($CURRENT_VERSION) est compatible avec macOS 26+"
    exit 0
fi

echo ""
echo "❌ Version de Go incompatible avec macOS 26+"
echo "   macOS 26+ nécessite Go 1.24+ pour générer le LC_UUID requis."
echo ""

if ! command -v brew &> /dev/null; then
    echo "⚠️  Homebrew n'est pas installé."
    echo "   Installez Go manuellement depuis: https://golang.org/dl/"
    echo "   Ou installez Homebrew: https://brew.sh/"
    exit 1
fi

echo "📦 Mise à jour de Go via Homebrew..."
echo ""

# Installer la dernière version de Go
if brew list go &> /dev/null; then
    echo "Mise à jour de Go..."
    brew upgrade go
else
    echo "Installation de Go..."
    brew install go
fi

echo ""
echo "✅ Mise à jour terminée!"
echo ""

# Vérifier quelle version est dans le PATH
CURRENT_GO=$(which go)
HOMEBREW_GO="/opt/homebrew/bin/go"

if [ -f "$HOMEBREW_GO" ]; then
    NEW_VERSION=$($HOMEBREW_GO version | awk '{print $3}')
    echo "Vérification de la nouvelle version:"
    $HOMEBREW_GO version
    echo ""
    
    if [ "$CURRENT_GO" != "$HOMEBREW_GO" ]; then
        echo "⚠️  ATTENTION: Homebrew Go ($NEW_VERSION) installé mais pas dans le PATH"
        echo "   Le PATH actuel pointe vers: $CURRENT_GO"
        echo ""
        echo "Pour utiliser Go 1.25.3, ajoutez ceci à votre ~/.zshrc ou ~/.bashrc:"
        echo "   export PATH=\"/opt/homebrew/bin:\$PATH\""
        echo ""
        echo "Ou utilisez directement:"
        echo "   /opt/homebrew/bin/go build -o comptes ./cmd/comptes"
        echo ""
        echo "Test avec la nouvelle version:"
        /opt/homebrew/bin/go build -o comptes ./cmd/comptes && echo "✓ Compilation réussie avec Go $NEW_VERSION"
    else
        echo "✅ Go $NEW_VERSION est maintenant dans votre PATH"
        echo "Vous pouvez compiler avec: make build"
    fi
else
    echo "Vérification de la version:"
    go version
    echo ""
    echo "Vous pouvez maintenant compiler avec: make build"
fi

