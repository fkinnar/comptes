# TODO List - Progression vers le MVP

## 📊 Progression actuelle : 95% vers le MVP (19/20)

---

## ✅ COMPLETÉ (19/20)

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

---

## ⏳ RESTANT À FAIRE (1/20)

### 🧪 Tests & Qualité
- ⏳ Tests d'intégration complets du MVP

---

## 🎯 Prochaines étapes prioritaires pour le MVP

1. **Tests d'intégration** complets
2. **Documentation** des commandes avancées
3. **Validation** des cas limites

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

---

## 📝 Notes

- **Architecture solide** : Base extensible prête pour l'évolution
- **Tests complets** : Couverture Service, Config, Storage
- **Infrastructure** : GitHub + pre-commit hooks
- **Documentation** : README, SETUP, architecture détaillée
- **Vision produit** : Fonctionnalités avancées identifiées pour l'évolution
- **Interface Git-like** : Audit trail complet avec undo intelligent
- **UUID courts** : Interface familière pour les développeurs
- **Formats multiples** : CSV pour Nushell, JSON pour scripting

---

*Dernière mise à jour : 28 octobre 2025*
