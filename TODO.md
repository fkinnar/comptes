# TODO List - Progression vers le MVP

## 📊 Progression actuelle : 100% vers le MVP (20/20)

---

## ✅ COMPLETÉ (20/20)

### 🏗️ Architecture & Infrastructure
- ✅ Architecture en couches (Service → Storage → Domain)
- ✅ Modèles de données (Account, Transaction, Category, Tag)
- ✅ Interface Storage abstraite pour persistance
- ✅ Implémentation JSONStorage
- ✅ Système de configuration YAML
- ✅ Tests unitaires complets (Service, Config, Storage)
- ✅ Repo GitHub avec pre-commit hooks

### ⚙️ Fonctionnalités de base
- ✅ Service de gestion des transactions (ajout, validation, calcul solde)
- ✅ CLI de base (init, add, list, balance)
- ✅ Support des dates flexibles dans JSON (today, yesterday, 2024-01-15)
- ✅ Renommage account_id → account pour JSON plus propre
- ✅ Interface JSON pure (suppression du flag --date)

### 🎯 CLI Avancé
- ✅ Commande edit (soft delete + nouvelle transaction)
- ✅ Commande delete (soft delete avec commentaire)
- ✅ Commande undo (détection automatique du type d'opération)
- ✅ Flag --history pour voir toutes les transactions
- ✅ Flag --format avec options text, csv, json
- ✅ Messages obligatoires pour edit/delete (-m, --message)
- ✅ UUID courts style Git (8 caractères)
- ✅ Support des IDs partiels (edit fd66)

### 🔧 Fonctionnalités Git-like
- ✅ Audit trail complet avec commentaires
- ✅ Relations parent-enfant pour les edits
- ✅ Undo intelligent (delete/add/edit)
- ✅ Historique complet des modifications
- ✅ Interface familière pour les développeurs

### 🧪 Tests & Qualité
- ✅ Tests d'intégration complets du MVP
- ✅ Plan de test complet avec edge cases
- ✅ Script de test automatique (28 tests)
- ✅ Pre-commit hooks avec validation complète
- ✅ Configuration par défaut pour initialisation

---

## 🎉 MVP COMPLET ! (20/20)

---

## 🎯 Prochaines étapes (post-MVP) - Priorités définies

### 🚀 Priorité 1 : Ergonomie quotidienne (CRUCIAL)
1. **Mode transactionnel avec contexte + Flags pour add** - Réduit drastiquement la verbosité
   ```bash
   # Mode transactionnel
   comptes account BANQUE
   comptes category ALM
   comptes add '{"amount": -25.50, "description": "Courses"}'
   comptes commit
   
   # Ou avec flags directement
   ./comptes add -a -25.50 -d "Courses" -c ALM -t URG
   ```
2. **Support multi-comptes avec transferts** - Gestion réaliste des finances
   ```bash
   comptes add '{"account": "BANQUE", "amount": -100, "transfer_to": "LIVRET"}'
   ```

### 🚀 Priorité 2 : Intégration pratique
3. **Import CSV** - Intégration avec relevés bancaires
   ```bash
   comptes add --file bank_statement.csv
   comptes add --csv "date,amount,description,category"
   ```

### 🚀 Priorité 3 : Analytics basiques
4. **Rapports simples** - Vision claire des finances
   ```bash
   comptes report --month 2024-01
   comptes report --category ALM --from 2024-01-01
   comptes balance --trend
   ```

### 🚀 Priorité 4 : Personnalisation
5. **Gestion catégories/tags via CLI** - Personnalisation sans fichiers
   ```bash
   comptes category add "VET" "Vêtements"
   comptes tag add "IMP" "Important"
   ```

---

## 🔧 Améliorations UX immédiates (à implémenter rapidement)

### 📝 Interface utilisateur améliorée
- ✅ **Aide contextuelle** : `--categories (-c)` et `--tags (-t)` sur `list` pour voir les options
  ```bash
  ./comptes list --categories  # Affiche toutes les catégories disponibles
  ./comptes list --tags        # Affiche tous les tags disponibles
  ```
- ✅ **Affichage amélioré** : Noms complets des catégories/tags au lieu des codes
  ```bash
  ./comptes list  # Affiche "Alimentation" au lieu de "ALM"
  ./comptes list --codes  # Flag pour garder les codes si besoin
  ```
- ✅ **Support des formats CSV/JSON** : Pour catégories, tags et transactions
  ```bash
  ./comptes list --categories --format csv  # Export CSV des catégories
  ./comptes list --tags --format json        # Export JSON des tags
  ./comptes list --transactions --format csv # Export CSV des transactions
  ```
- ✅ **Architecture cohérente** : Flag `--transactions` par défaut pour clarté
- ✅ **CSV compatible Nushell** : Échappement correct des virgules dans les descriptions
- **Flags pour add** : `--amount (-a)`, `--description (-d)`, `--categories (-c)`, `--tags (-t)`, `--date` (implémenté avec le mode transactionnel)

### 🗑️ Opérations avancées
- **Suppression définitive** : `--hard` pour `delete`, `edit`, `undo`
  ```bash
  ./comptes delete abc123 --hard -m "Dupliquée"  # Suppression définitive
  ./comptes edit abc123 '{"amount": 100}' --hard -m "Correction définitive"
  ./comptes undo def456 --hard  # Suppression définitive au lieu de désactivation
  ```
- **Confirmation forcée** : `-f/--force` pour bypasser les confirmations
  ```bash
  ./comptes delete abc123 --hard --force -m "Dupliquée"  # Pas de confirmation
  ```

### 🎨 Améliorations d'affichage
- **Statut intelligent** : Pas de ✅/❌ pour `list` normal (toutes sont actives)
- **Historique clair** : ✅/❌ seulement pour `list --history`
- **Messages informatifs** : "removed" au lieu de "deactivated" pour undo edit

---

## ⚡ Améliorations techniques importantes (à faire rapidement)

### 🗄️ Performance et scalabilité
- **Snapshots de solde** : Éviter de recalculer depuis le début
  ```go
  // Ajouter des snapshots périodiques pour accélérer les calculs
  type BalanceSnapshot struct {
      AccountID string    `json:"account_id"`
      Balance   float64   `json:"balance"`
      Date      time.Time `json:"date"`
      LastTxnID string    `json:"last_transaction_id"`
  }
  ```
- **Couche de cache** : Entre Service et Storage pour optimiser les lectures
- **Indexation** : Pour les recherches rapides par date, catégorie, etc.

### 🗃️ Storage avancé
- **Migration SQLite** : Quand la structure sera validée à l'usage
  ```go
  // Nouvelle implémentation Storage
  type SQLiteStorage struct {
      db *sql.DB
  }
  ```
- **Couche d'abstraction** : Au-dessus du Storage pour la logique métier
- **Migration automatique** : JSON → SQLite sans perte de données

### 🔄 Architecture évolutive
- **Interface Storage enrichie** : Méthodes pour snapshots, indexation
- **Service layer étendu** : Cache, validation avancée, analytics
- **Configuration dynamique** : Modification des catégories/tags sans redémarrage

---

## 🧪 Approche de validation et amélioration continue

### 📊 Utilisation pour validation
- **Usage quotidien** : Identifier les points de friction
- **Tests de charge** : Avec de vraies données (1000+ transactions)
- **Validation structure** : Confirmer que les modèles sont adaptés
- **Feedback utilisateur** : Améliorer l'ergonomie basée sur l'usage réel

### 🔄 Cycle d'amélioration
1. **Utiliser l'outil** avec des données réelles
2. **Identifier les problèmes** de performance/ergonomie
3. **Implémenter les corrections** prioritaires
4. **Valider les améliorations** avec de nouveaux tests
5. **Répéter** jusqu'à satisfaction

### 🎯 Objectif : Outil fonctionnel en quelques jours
- **Priorités 1-2** : Ergonomie quotidienne (mode transactionnel + multi-comptes)
- **Priorité 3** : Import CSV pour intégration
- **Priorité 4** : Analytics basiques pour vision
- **Améliorations techniques** : Snapshots, SQLite, cache

---

## 🚀 Fonctionnalités avancées (post-MVP)

### 🔄 Go routines & asynchrone
- **Analytics** : Calculs en parallèle sur plusieurs comptes
- **Import/Export** : Traitement de gros fichiers CSV/JSON
- **Validation** : Vérification des catégories/tags en parallèle
- **Cache** : Mise à jour asynchrone des soldes

### 🎭 Mode transactionnel avec contexte
- **Contexte partagé** : `comptes account BANQUE` → `comptes category ALM`
- **Transaction atomique** : `comptes commit` (tout ou rien)
- **Moins verbeux** : Plus besoin de JSON pour chaque transaction
- **Validation groupée** : Vérification à la fin

### 📅 Gestion des dates
- ✅ **Support des dates flexibles dans JSON** : `today`, `yesterday`, `2024-01-15`
- ✅ **Interface JSON pure** : Suppression du flag `--date` pour éviter les conflits
- ✅ **Parser intelligent** : FlexibleDate type avec UnmarshalJSON personnalisé
- ⏳ **Formats étendus** : `last week`, `next month`, etc. (post-MVP)

### 🧮 Calculs dans les requêtes
- **Expressions** : `{45.00 - 12.00}` pour les calculs
- **Références** : Montants dynamiques entre transactions
- **Validation** : Vérification des calculs et références

### 💼 Fonctionnalités métier avancées
- **Support complet multi-comptes** avec transferts
- **Gestion des catégories** (CRUD via CLI)
- **Gestion des tags** (CRUD via CLI)
- **Règles de validation avancées** (catégories existantes, etc.)

---

## 🧪 Guide de test complet

### 📋 Commandes de base à tester
```bash
# Initialisation
./comptes init

# Ajout de transactions
./comptes add '{"account": "BANQUE", "amount": -25.50, "description": "Test", "categories": ["ALM"]}'
./comptes add '{"account": "BANQUE", "amount": 1500, "description": "Salaire", "categories": ["SLR"], "date": "today"}'

# Liste avec différents formats
./comptes list
./comptes list --history
./comptes list --format csv
./comptes list --format json
./comptes list --history --format csv

# Nouvelles fonctionnalités d'aide contextuelle
./comptes list --categories              # Affiche les catégories disponibles
./comptes list --categories --format csv # Export CSV des catégories
./comptes list --tags --format json     # Export JSON des tags
./comptes list --codes                  # Affiche les codes au lieu des noms

# Édition avec message obligatoire
./comptes edit <id> '{"amount": -30.00}' -m "Correction montant"
./comptes edit <id> '{"description": "Nouvelle description"}' --message "Fix typo"

# Suppression avec message obligatoire
./comptes delete <id> -m "Transaction erronée"

# Undo intelligent
./comptes undo <id>  # Détecte automatiquement le type d'opération

# Solde
./comptes balance

# Migration des IDs
./comptes migrate
```

### 🔍 Cas de test spécifiques
1. **IDs partiels** : `edit fd66` au lieu de `edit fd6647d8`
2. **Dates flexibles** : `"date": "today"`, `"date": "yesterday"`
3. **Messages obligatoires** : Tester sans `-m` pour edit/delete
4. **Undo** : Tester undo add, undo edit, undo delete
5. **Historique** : Vérifier que `--history` montre les commentaires
6. **Formats** : Vérifier CSV pour Nushell, JSON pour scripting
7. **Aide contextuelle** : Tester `--categories` et `--tags` avec différents formats
8. **Affichage amélioré** : Vérifier noms complets vs codes avec `--codes`
9. **CSV Nushell** : Tester `./comptes list --categories --format csv | from csv`
10. **Architecture cohérente** : Vérifier que `--transactions` fonctionne comme par défaut

---

## 📝 Notes

- **Architecture solide** : Base extensible prête pour l'évolution
- **Tests complets** : Couverture Service, Config, Storage + tests d'intégration
- **Infrastructure** : GitHub + pre-commit hooks avec validation complète
- **Documentation** : README, SETUP, architecture détaillée, plan de test
- **Vision produit** : Fonctionnalités avancées identifiées pour l'évolution
- **Interface Git-like** : Audit trail complet avec undo intelligent
- **UUID courts** : Interface familière pour les développeurs
- **Formats multiples** : CSV pour Nushell, JSON pour scripting
- **Configuration par défaut** : Initialisation automatique sans fichiers
- **Suite de tests** : 28 tests automatiques avec edge cases

---

*Dernière mise à jour : 28 octobre 2025 - Nouvelles fonctionnalités UX implémentées*
