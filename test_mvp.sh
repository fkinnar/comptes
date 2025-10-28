#!/bin/bash
# test_mvp.sh - Script de test automatique pour le MVP

set -e  # Arrêter en cas d'erreur

echo "🧪 Début des tests MVP Comptes..."
echo "=================================="

# Couleurs pour les logs
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Compteurs
TESTS_PASSED=0
TESTS_FAILED=0
TOTAL_TESTS=0

# Fonction pour exécuter un test
run_test() {
    local test_name="$1"
    local test_command="$2"
    local expected_exit_code="${3:-0}"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -n "Test $TOTAL_TESTS: $test_name... "
    
    if eval "$test_command" >/dev/null 2>&1; then
        if [ $? -eq $expected_exit_code ]; then
            echo -e "${GREEN}✅ PASS${NC}"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            echo -e "${RED}❌ FAIL (exit code: $?)${NC}"
            TESTS_FAILED=$((TESTS_FAILED + 1))
        fi
    else
        if [ $? -eq $expected_exit_code ]; then
            echo -e "${GREEN}✅ PASS${NC}"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            echo -e "${RED}❌ FAIL (exit code: $?)${NC}"
            TESTS_FAILED=$((TESTS_FAILED + 1))
        fi
    fi
}

# Fonction pour vérifier la sortie
check_output() {
    local test_name="$1"
    local command="$2"
    local expected_pattern="$3"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    echo -n "Test $TOTAL_TESTS: $test_name... "
    
    if eval "$command" | grep -q "$expected_pattern"; then
        echo -e "${GREEN}✅ PASS${NC}"
        TESTS_PASSED=$((TESTS_PASSED + 1))
    else
        echo -e "${RED}❌ FAIL (pattern not found)${NC}"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
}

# Nettoyage initial
echo "🧹 Nettoyage initial..."
rm -rf test-data/ test-config/ 2>/dev/null || true

# Créer les dossiers de test
mkdir -p test-data test-config

# Variables d'environnement pour les tests
export COMPTES_DATA_DIR="test-data"
export COMPTES_CONFIG_DIR="test-config"

# Compilation
echo "🔨 Compilation..."
go build -o comptes cmd/comptes/main.go

echo ""
echo "📋 Tests de base..."
echo "=================="

# Test 1: Initialisation
run_test "Initialisation propre" "./comptes init" 0
run_test "Réinitialisation" "./comptes init" 0

# Vérifier que les fichiers sont créés
if [ ! -f "test-config/config.yaml" ] || [ ! -f "test-data/transactions.json" ]; then
    echo -e "${RED}❌ Fichiers de configuration non créés${NC}"
    exit 1
fi

# Test 2: Ajout de transactions
run_test "Transaction simple" "./comptes add '{\"account\": \"BANQUE\", \"amount\": -25.50, \"description\": \"Test\", \"categories\": [\"ALM\"]}'" 0
run_test "Transaction avec date relative" "./comptes add '{\"account\": \"BANQUE\", \"amount\": 1500, \"description\": \"Salaire\", \"categories\": [\"SLR\"], \"date\": \"today\"}'" 0
run_test "Transaction avec tags" "./comptes add '{\"account\": \"BANQUE\", \"amount\": -15, \"description\": \"Test tags\", \"categories\": [\"ALM\"], \"tags\": [\"URG\", \"REC\"]}'" 0

# Test 3: Liste et formats
run_test "Liste normale" "./comptes list" 0
run_test "Liste avec historique" "./comptes list --history" 0
run_test "Format CSV" "./comptes list --format csv" 0
run_test "Format JSON" "./comptes list --format json" 0
run_test "Format CSV avec historique" "./comptes list --history --format csv" 0

# Test 4: Vérifications de sortie
check_output "Format CSV valide" "./comptes list --format csv" "id,date,amount,description,categories,tags"
check_output "Format JSON valide" "./comptes list --format json" "\["
check_output "Indicateurs d'état" "./comptes list --history" "✅\|❌"

echo ""
echo "🚨 Tests d'erreurs..."
echo "==================="

# Test 5: Erreurs de validation
run_test "JSON invalide" "./comptes add '{\"account\": \"BANQUE\", \"amount\": -25.50, \"description\": \"Test\", \"categories\": [\"ALM\"]'" 1
run_test "Champs manquants" "./comptes add '{\"amount\": -25.50, \"description\": \"Test\"}'" 1
run_test "Catégorie inexistante" "./comptes add '{\"account\": \"BANQUE\", \"amount\": -25.50, \"description\": \"Test\", \"categories\": [\"INEXISTANT\"]}'" 1

# Test 6: Erreurs de commande
run_test "Commande inexistante" "./comptes unknown_command" 1
run_test "Arguments manquants pour edit" "./comptes edit" 1
run_test "Message manquant pour edit" "./comptes edit nonexistent '{\"amount\": -30.00}'" 1
run_test "Message manquant pour delete" "./comptes delete nonexistent" 1

# Test 7: IDs et références
run_test "ID inexistant" "./comptes edit nonexistent '{\"amount\": -30.00}' -m \"Test\"" 1

# Test 7.5: Nouvelles fonctionnalités de list
echo ""
echo "📋 Test des nouvelles fonctionnalités de list"
echo "============================================="

# Ajouter quelques transactions avec catégories et tags pour tester
./comptes add '{"account": "BANQUE", "amount": -25.50, "description": "Test categories", "categories": ["ALM"], "tags": ["REC"]}' >/dev/null 2>&1
./comptes add '{"account": "BANQUE", "amount": 2800, "description": "Test tags", "categories": ["SLR"], "tags": ["URG"]}' >/dev/null 2>&1

# Test aide contextuelle
run_test "Aide contextuelle - catégories" "./comptes list --categories" 0
run_test "Aide contextuelle - tags" "./comptes list --tags" 0
run_test "Aide contextuelle - catégories court" "./comptes list -c" 0
run_test "Aide contextuelle - tags court" "./comptes list -t" 0

# Test formats CSV et JSON pour catégories et tags
run_test "Catégories en CSV" "./comptes list --categories --format csv" 0
run_test "Catégories en JSON" "./comptes list --categories --format json" 0
run_test "Tags en CSV" "./comptes list --tags --format csv" 0
run_test "Tags en JSON" "./comptes list --tags --format json" 0

# Test affichage avec noms complets vs codes
run_test "Liste avec noms complets" "./comptes list" 0
run_test "Liste avec codes" "./comptes list --codes" 0

# Test formats pour transactions
run_test "Transactions en CSV (noms complets)" "./comptes list --transactions --format csv" 0
run_test "Transactions en CSV (codes)" "./comptes list --transactions --format csv --codes" 0
run_test "Transactions en JSON" "./comptes list --transactions --format json" 0

# Test flag --transactions explicite
run_test "Transactions explicites" "./comptes list --transactions" 0

# Test combinaisons complexes
run_test "Liste avec history et codes" "./comptes list --history --codes" 0
run_test "Liste avec history, format csv et codes" "./comptes list --history --format csv --codes" 0
run_test "Liste avec transactions, history et format json" "./comptes list --transactions --history --format json" 0

# Test cas d'erreur pour formats invalides
run_test "Format invalide (fallback vers text)" "./comptes list --categories --format invalid" 0

# Test flags conflictuels (categories vs tags - categories gagne)
run_test "Flags conflictuels (categories prioritaire)" "./comptes list --categories --tags" 0

# Test cas d'erreur edge
run_test "Catégorie inexistante" "./comptes add '{\"account\": \"BANQUE\", \"amount\": -50, \"description\": \"Test\", \"categories\": [\"INEXISTANT\"]}'" 1
run_test "Tag inexistant" "./comptes add '{\"account\": \"BANQUE\", \"amount\": -50, \"description\": \"Test\", \"tags\": [\"INEXISTANT\"]}'" 1

echo ""
echo "🔄 Tests de cohérence..."
echo "======================="

# Test 8: Édition et audit trail
run_test "Édition simple" "./comptes edit \$(./comptes list --format csv | tail -n +2 | head -1 | cut -d',' -f1) '{\"amount\": -30.00}' -m \"Correction montant\"" 0
run_test "Suppression" "./comptes delete \$(./comptes list --format csv | tail -n +2 | head -1 | cut -d',' -f1) -m \"Transaction erronée\"" 0

# Test 9: Undo
run_test "Undo delete" "./comptes undo \$(./comptes list --history --format csv | grep 'false' | head -1 | cut -d',' -f1)" 0

# Test 10: Solde
run_test "Calcul de solde" "./comptes balance" 0

echo ""
echo "🧪 Tests de performance..."
echo "========================"

# Test 11: Nombreuses transactions
echo "Ajout de 50 transactions..."
for i in {1..50}; do
    ./comptes add "{\"account\": \"BANQUE\", \"amount\": -$i, \"description\": \"Test $i\", \"categories\": [\"ALM\"]}" >/dev/null 2>&1
done

run_test "Liste avec 50+ transactions" "./comptes list" 0
run_test "Calcul de solde avec 50+ transactions" "./comptes balance" 0

# Test 12: Nouvelles fonctionnalités --hard et --force
echo ""
echo "🔧 Test des nouvelles fonctionnalités --hard et --force"
echo "======================================================"

# Ajouter une transaction pour tester
./comptes add '{"account": "BANQUE", "amount": -100, "description": "Test hard delete", "categories": ["ALM"]}' >/dev/null 2>&1
TEST_TXN_ID=$(./comptes list --format csv | tail -n +2 | head -1 | cut -d',' -f1)

# Test suppression définitive avec confirmation (simuler "y")
echo "y" | run_test "Suppression définitive avec confirmation" "./comptes delete $TEST_TXN_ID --hard -m \"Test suppression définitive\"" 0

# Ajouter une autre transaction pour tester --force
./comptes add '{"account": "BANQUE", "amount": -200, "description": "Test force delete", "categories": ["ALM"]}' >/dev/null 2>&1
TEST_TXN_ID=$(./comptes list --format csv | tail -n +2 | head -1 | cut -d',' -f1)

# Test suppression définitive avec --force (pas de confirmation)
run_test "Suppression définitive avec --force" "./comptes delete $TEST_TXN_ID --hard --force -m \"Test suppression définitive avec force\"" 0

# Ajouter une transaction pour tester undo --hard
./comptes add '{"account": "BANQUE", "amount": -300, "description": "Test undo hard", "categories": ["ALM"]}' >/dev/null 2>&1
TEST_TXN_ID=$(./comptes list --format csv | tail -n +2 | head -1 | cut -d',' -f1)

# Test undo définitif avec --force
run_test "Undo définitif avec --force" "./comptes undo $TEST_TXN_ID --hard --force" 0

# Test annulation de confirmation (simuler "n")
echo "n" | run_test "Annulation confirmation delete --hard" "./comptes delete test123 --hard -m \"Test annulation\"" 0
echo "n" | run_test "Annulation confirmation undo --hard" "./comptes undo test456 --hard" 0

# Test 13: Migration
run_test "Migration des IDs" "./comptes migrate" 0

echo ""
echo "📊 Résultats des tests"
echo "====================="
echo -e "Tests passés: ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests échoués: ${RED}$TESTS_FAILED${NC}"
echo -e "Total: $TOTAL_TESTS"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "\n${GREEN}🎉 Tous les tests sont passés !${NC}"
    exit 0
else
    echo -e "\n${RED}❌ $TESTS_FAILED test(s) ont échoué${NC}"
    exit 1
fi
