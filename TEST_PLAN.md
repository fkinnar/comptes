# Plan de test complet - Comptes MVP

## 🎯 Objectif
Valider que toutes les fonctionnalités du MVP fonctionnent correctement ensemble, y compris les cas limites et les erreurs.

---

## 📋 Tests de base (Happy Path)

### 1. Initialisation
```bash
# Test 1.1: Initialisation propre
rm -rf data/ config/
./comptes init
# Vérifier: data/ et config/ créés, fichiers JSON vides

# Test 1.2: Réinitialisation
./comptes init
# Vérifier: Pas d'erreur, fichiers préservés
```

### 2. Ajout de transactions
```bash
# Test 2.1: Transaction simple
./comptes add '{"account": "BANQUE", "amount": -25.50, "description": "Test", "categories": ["ALM"]}'
# Vérifier: Transaction ajoutée, ID généré

# Test 2.2: Transaction avec date relative
./comptes add '{"account": "BANQUE", "amount": 1500, "description": "Salaire", "categories": ["SLR"], "date": "today"}'
# Vérifier: Date = aujourd'hui

# Test 2.3: Transaction avec date absolue
./comptes add '{"account": "BANQUE", "amount": -100, "description": "Loyer", "categories": ["LGT"], "date": "2024-01-15"}'
# Vérifier: Date = 2024-01-15

# Test 2.4: Transaction avec tags
./comptes add '{"account": "BANQUE", "amount": -15, "description": "Test tags", "categories": ["ALM"], "tags": ["URG", "REC"]}'
# Vérifier: Tags affichés dans list
```

### 3. Liste et formats
```bash
# Test 3.1: Liste normale
./comptes list
# Vérifier: Seulement transactions actives

# Test 3.2: Liste avec historique
./comptes list --history
# Vérifier: Toutes les transactions, indicateurs ✅/❌

# Test 3.3: Format CSV
./comptes list --format csv
# Vérifier: En-têtes corrects, données séparées par virgules

# Test 3.4: Format JSON
./comptes list --format json
# Vérifier: JSON valide, structure correcte

# Test 3.5: Format CSV avec historique
./comptes list --history --format csv
# Vérifier: Colonne is_active présente
```

### 4. Édition
```bash
# Test 4.1: Édition simple
./comptes edit <id> '{"amount": -30.00}' -m "Correction montant"
# Vérifier: Ancienne transaction ❌, nouvelle ✅

# Test 4.2: Édition avec ID partiel
./comptes edit <id_partiel> '{"description": "Nouvelle description"}' --message "Fix typo"
# Vérifier: Fonctionne avec 4+ caractères

# Test 4.3: Édition multiple champs
./comptes edit <id> '{"amount": -40.00, "description": "Modifié", "categories": ["LGT"]}' -m "Changement complet"
# Vérifier: Tous les champs modifiés
```

### 5. Suppression
```bash
# Test 5.1: Suppression simple
./comptes delete <id> -m "Transaction erronée"
# Vérifier: Transaction ❌ avec commentaire

# Test 5.2: Suppression avec ID partiel
./comptes delete <id_partiel> --message "Suppression test"
# Vérifier: Fonctionne avec ID partiel
```

### 6. Undo
```bash
# Test 6.1: Undo delete
./comptes undo <id_supprimé>
# Vérifier: Transaction restaurée ✅

# Test 6.2: Undo add
./comptes undo <id_ajouté>
# Vérifier: Transaction désactivée ❌

# Test 6.3: Undo edit
./comptes undo <id_enfant>
# Vérifier: Parent restauré ✅, enfant désactivé ❌
```

### 7. Solde
```bash
# Test 7.1: Calcul de solde
./comptes balance
# Vérifier: Solde = initial_balance + somme des transactions actives
```

---

## 🚨 Tests d'erreurs et cas limites

### 8. Erreurs de validation
```bash
# Test 8.1: JSON invalide
./comptes add '{"account": "BANQUE", "amount": -25.50, "description": "Test", "categories": ["ALM"]'
# Vérifier: Erreur de parsing JSON

# Test 8.2: Champs manquants
./comptes add '{"amount": -25.50, "description": "Test"}'
# Vérifier: Erreur de validation (account manquant)

# Test 8.3: Catégorie inexistante
./comptes add '{"account": "BANQUE", "amount": -25.50, "description": "Test", "categories": ["INEXISTANT"]}'
# Vérifier: Erreur de validation

# Test 8.4: Montant zéro
./comptes add '{"account": "BANQUE", "amount": 0, "description": "Test", "categories": ["ALM"]}'
# Vérifier: Erreur ou transaction créée (selon validation)
```

### 9. Erreurs de commande
```bash
# Test 9.1: Commande inexistante
./comptes unknown_command
# Vérifier: Message d'erreur avec usage

# Test 9.2: Arguments manquants
./comptes edit
# Vérifier: Usage avec exemples

# Test 9.3: Message manquant pour edit
./comptes edit <id> '{"amount": -30.00}'
# Vérifier: Erreur "Message is mandatory"

# Test 9.4: Message manquant pour delete
./comptes delete <id>
# Vérifier: Erreur "Message is mandatory"
```

### 10. IDs et références
```bash
# Test 10.1: ID inexistant
./comptes edit nonexistent '{"amount": -30.00}' -m "Test"
# Vérifier: Erreur "no transaction found"

# Test 10.2: ID partiel ambigu
# Créer deux transactions avec IDs commençant par les mêmes caractères
./comptes edit ab '{"amount": -30.00}' -m "Test"
# Vérifier: Erreur "multiple transactions found"

# Test 10.3: Suppression d'une transaction déjà supprimée
./comptes delete <id_déjà_supprimé> -m "Test"
# Vérifier: Erreur "already deleted"
```

### 11. Dates et formats
```bash
# Test 11.1: Date invalide
./comptes add '{"account": "BANQUE", "amount": -25.50, "description": "Test", "categories": ["ALM"], "date": "invalid-date"}'
# Vérifier: Erreur de parsing de date

# Test 11.2: Format invalide
./comptes list --format invalid
# Vérifier: Fallback vers format text

# Test 11.3: Dates relatives complexes
./comptes add '{"account": "BANQUE", "amount": -25.50, "description": "Test", "categories": ["ALM"], "date": "yesterday"}'
# Vérifier: Date = hier
```

### 12. Fichiers et permissions
```bash
# Test 12.1: Répertoire data/ supprimé
rm -rf data/
./comptes list
# Vérifier: Erreur ou création automatique

# Test 12.2: Fichier JSON corrompu
echo "invalid json" > data/transactions.json
./comptes list
# Vérifier: Erreur de parsing ou récupération

# Test 12.3: Permissions en lecture seule
chmod 444 data/transactions.json
./comptes add '{"account": "BANQUE", "amount": -25.50, "description": "Test", "categories": ["ALM"]}'
# Vérifier: Erreur de permission
```

---

## 🔄 Tests de cohérence et intégrité

### 13. Audit trail
```bash
# Test 13.1: Chaîne d'éditions
./comptes add '{"account": "BANQUE", "amount": -25.50, "description": "Test", "categories": ["ALM"]}'
./comptes edit <id1> '{"amount": -30.00}' -m "Edit 1"
./comptes edit <id2> '{"description": "Modifié"}' -m "Edit 2"
./comptes list --history
# Vérifier: 3 transactions, relations parent-enfant correctes

# Test 13.2: Undo en chaîne
./comptes undo <id3>  # Undo edit 2
./comptes undo <id2>  # Undo edit 1
./comptes list --history
# Vérifier: Retour à l'état initial
```

### 14. Performance et limites
```bash
# Test 14.1: Nombreuses transactions
for i in {1..100}; do
  ./comptes add "{\"account\": \"BANQUE\", \"amount\": -$i, \"description\": \"Test $i\", \"categories\": [\"ALM\"]}"
done
./comptes list --format csv | wc -l
# Vérifier: 100 transactions + en-tête

# Test 14.2: Calcul de solde avec beaucoup de transactions
./comptes balance
# Vérifier: Calcul correct et rapide
```

### 15. Migration et compatibilité
```bash
# Test 15.1: Migration des anciens IDs
# Créer manuellement des transactions avec anciens IDs
echo '[{"id": "txn_1234567890", "account": "BANQUE", "amount": -25.50, "description": "Test", "categories": ["ALM"], "is_active": true, "created_at": "2024-01-15T10:30:00Z", "updated_at": "2024-01-15T10:30:00Z"}]' > data/transactions.json
./comptes migrate
./comptes list
# Vérifier: IDs migrés vers UUID courts
```

---

## 🧪 Tests d'intégration avec outils externes

### 16. Intégration Nushell
```bash
# Test 16.1: Import CSV
./comptes list --format csv | head -5
# Vérifier: Format CSV valide

# Test 16.2: Filtrage
./comptes list --history --format csv | grep "false"
# Vérifier: Transactions supprimées visibles
```

### 17. Intégration JSON
```bash
# Test 17.1: Parsing JSON
./comptes list --format json | jq '.[0].amount'
# Vérifier: Montant correct

# Test 17.2: Filtrage JSON
./comptes list --history --format json | jq '.[] | select(.is_active == false)'
# Vérifier: Transactions supprimées
```

---

## 📊 Critères de réussite

### ✅ Tests de base (100% requis)
- [ ] Initialisation fonctionne
- [ ] Ajout de transactions fonctionne
- [ ] Liste avec tous les formats fonctionne
- [ ] Édition avec messages fonctionne
- [ ] Suppression avec messages fonctionne
- [ ] Undo intelligent fonctionne
- [ ] Calcul de solde fonctionne

### ✅ Tests d'erreurs (90% requis)
- [ ] Gestion des erreurs de validation
- [ ] Gestion des erreurs de commande
- [ ] Gestion des IDs inexistants/ambigus
- [ ] Gestion des dates invalides
- [ ] Gestion des fichiers corrompus

### ✅ Tests de cohérence (95% requis)
- [ ] Audit trail complet
- [ ] Relations parent-enfant correctes
- [ ] Undo en chaîne fonctionne
- [ ] Migration des IDs fonctionne

### ✅ Tests d'intégration (80% requis)
- [ ] Format CSV compatible Nushell
- [ ] Format JSON compatible jq
- [ ] Performance acceptable (100+ transactions)

---

## 🚀 Exécution des tests

### Script de test automatique
```bash
#!/bin/bash
# test_mvp.sh

echo "🧪 Début des tests MVP..."

# Tests de base
echo "📋 Tests de base..."
# ... (implémentation des tests)

# Tests d'erreurs
echo "🚨 Tests d'erreurs..."
# ... (implémentation des tests)

# Tests de cohérence
echo "🔄 Tests de cohérence..."
# ... (implémentation des tests)

echo "✅ Tests terminés!"
```

### Intégration pre-commit
- Exécution automatique des tests critiques
- Validation des formats de sortie
- Vérification de la cohérence des données
