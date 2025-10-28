# TODO List - Progression vers le MVP

## 📊 Progression actuelle : 50% vers le MVP (9/18)

---

## ✅ COMPLETÉ (9/18)

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

---

## 🔄 EN COURS (0/18)

---

## ⏳ RESTANT À FAIRE (9/18)

### 🔧 CLI Avancé
- ⏳ Commande edit (soft delete + nouvelle transaction)
- ⏳ Commande delete (soft delete avec commentaire)
- ⏳ Flag --history pour voir toutes les transactions
- ⏳ Support JSON pour add/edit (--json, --file)

### 💼 Fonctionnalités métier
- ⏳ Support complet multi-comptes avec transferts
- ⏳ Gestion des catégories (CRUD via CLI)
- ⏳ Gestion des tags (CRUD via CLI)
- ⏳ Règles de validation avancées (catégories existantes, etc.)

### 🧪 Tests & Qualité
- ⏳ Tests d'intégration complets du MVP

---

## 🎯 Prochaines étapes prioritaires pour le MVP

1. **CLI edit/delete** (soft delete)
2. **Support multi-comptes** (transferts)
3. **Gestion des catégories/tags**
4. **Tests d'intégration**

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
- **Flag --date** : `--date 2024-01-15` ou `-d 2024-01-15`
- **Formats flexibles** : `2024-01-15`, `15/01/2024`, `yesterday`, `last week`
- **Validation** : Dates cohérentes et réalistes

### 🧮 Calculs dans les requêtes
- **Expressions** : `{45.00 - 12.00}` pour les calculs
- **Références** : Montants dynamiques entre transactions
- **Validation** : Vérification des calculs et références

---

## 📝 Notes

- **Architecture solide** : Base extensible prête pour l'évolution
- **Tests complets** : Couverture Service, Config, Storage
- **Infrastructure** : GitHub + pre-commit hooks
- **Documentation** : README, SETUP, architecture détaillée
- **Vision produit** : Fonctionnalités avancées identifiées pour l'évolution

---

*Dernière mise à jour : 27 octobre 2024*
