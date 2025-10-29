#!/bin/bash
# Script wrapper pour comptes - solution temporaire pour Go < 1.24
# Utilise go run au lieu d'un binaire compilé

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MAIN_PACKAGE="./cmd/comptes"

# Vérifier si Go est installé
if ! command -v go &> /dev/null; then
    echo "Erreur: Go n'est pas installé"
    exit 1
fi

# Vérifier la version de Go
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
GO_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
GO_MINOR=$(echo $GO_VERSION | cut -d. -f2)

# Avertir si version < 1.24
if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 24 ]); then
    echo "⚠️  Attention: Go $GO_VERSION détecté. macOS 26+ nécessite Go 1.24+ pour éviter l'erreur LC_UUID."
    echo "   Utilisation de 'go run' comme solution temporaire..."
    echo "   Mettez à jour Go avec: brew install go@1.25"
    echo ""
fi

# Exécuter avec go run
cd "$SCRIPT_DIR"
go run "$MAIN_PACKAGE" "$@"

