# TODO List - Progression du projet Comptes

## 📊 Progression actuelle : MVP + Améliorations UX complètes

---

## ✅ COMPLETÉ - MVP (20/20) + Améliorations UX (15/15)

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

### 🎨 Améliorations UX récentes
- ✅ **Aide contextuelle** : `--categories (-c)` et `--tags (-t)` sur `list`
- ✅ **Affichage amélioré** : Noms complets des catégories/tags au lieu des codes
- ✅ **Flag --codes** : Pour revenir aux codes si nécessaire
- ✅ **Support des formats CSV/JSON** : Pour catégories, tags et transactions
- ✅ **Architecture cohérente** : Flag `--transactions` par défaut pour clarté
- ✅ **CSV compatible Nushell** : Échappement correct des virgules dans les descriptions
- ✅ **Suppression définitive** : `--hard` pour `delete` et `undo`
- ✅ **Confirmation forcée** : `-f/--force` pour bypasser les confirmations
- ✅ **Aide mise à jour** : Documentation complète des nouveaux flags
- ✅ **Sécurité** : Confirmations obligatoires pour les opérations destructives
- ✅ **Flexibilité** : Possibilité de bypasser les confirmations avec `--force`

### 🧪 Tests & Qualité
- ✅ Tests d'intégration complets du MVP (52 tests)
- ✅ Tests unitaires complets (14 tests)
- ✅ Plan de test complet avec edge cases
- ✅ Script de test automatique avec nouvelles fonctionnalités
- ✅ Pre-commit hooks avec validation complète
- ✅ Configuration par défaut pour initialisation
- ✅ Couverture 100% des nouvelles fonctionnalités
- ✅ Tests des combinaisons complexes de flags
- ✅ Tests des cas d'erreur et edge cases
- ✅ Validation de compatibilité Nushell

---

## 🚀 PROCHAINES ÉTAPES - Priorités définies

### 🎯 Priorité 1 : Ergonomie quotidienne (CRUCIAL)
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

### 🎯 Priorité 2 : Intégration pratique
3. **Import CSV** - Intégration avec relevés bancaires
   ```bash
   comptes add --file bank_statement.csv
   comptes add --csv "date,amount,description,category"
   ```

### 🎯 Priorité 3 : Analytics basiques
4. **Rapports simples** - Vision claire des finances
   ```bash
   comptes report --month 2024-01
   comptes report --category ALM --from 2024-01-01
   comptes balance --trend
   ```

### 🎯 Priorité 4 : Personnalisation
5. **Gestion catégories/tags via CLI** - Personnalisation sans fichiers
   ```bash
   comptes category add "VET" "Vêtements"
   comptes tag add "IMP" "Important"
   ```

### 🎯 Priorité 5 : Déploiement production (CRUCIAL)
6. **Gestion des branches Git** - Séparation dev/production
   ```bash
   git checkout -b dev
   # main = production stable
   # dev = développement actif
   ```
7. **Installation système** - Fichiers aux bons endroits
   ```bash
   # Binaire dans le PATH
   sudo cp comptes /usr/local/bin/
   
   # Config et données dans les répertoires système
   mkdir -p ~/.config/comptes
   mkdir -p ~/.local/share/comptes
   
   # Ou respecter XDG Base Directory
   # ~/.config/comptes/config.yaml
   # ~/.local/share/comptes/movements.json
   ```

---

## 🔧 Améliorations techniques importantes (à faire rapidement)

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

### 🚀 Déploiement et distribution
- **Respect XDG Base Directory** : `~/.config/comptes/` et `~/.local/share/comptes/`
- **Installation système** : Binaire dans `/usr/local/bin/` ou `/opt/comptes/`
- **Gestion des branches** : `main` (stable) vs `dev` (développement)
- **Packaging** : Scripts d'installation pour différents OS
- **Variables d'environnement** : `COMPTES_CONFIG_DIR`, `COMPTES_DATA_DIR`

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

# Suppression définitive avec confirmation
./comptes delete <id> --hard -m "Dupliquée"

# Suppression définitive sans confirmation
./comptes delete <id> --hard --force -m "Dupliquée"

# Undo intelligent
./comptes undo <id>  # Détecte automatiquement le type d'opération

# Undo définitif
./comptes undo <id> --hard --force

# Solde
./comptes balance

# Migration des IDs (commande supprimée, plus nécessaire)
# ./comptes migrate
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
11. **Suppression définitive** : Tester `--hard` avec et sans `--force`
12. **Confirmations** : Tester les annulations (répondre "n")
13. **Combinaisons complexes** : Tester plusieurs flags ensemble
14. **Cas d'erreur** : Catégories/tags inexistants, formats invalides

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
- **Suite de tests** : 52 tests automatiques avec edge cases
- **UX moderne** : Aide contextuelle, affichage amélioré, confirmations intelligentes

---

*Dernière mise à jour : 28 octobre 2025 - MVP + Améliorations UX complètes*