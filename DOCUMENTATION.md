# Documentation Comptes - MVP Complet

## 📋 Table des matières

1. [Vue d'ensemble](#vue-densemble)
2. [Installation et configuration](#installation-et-configuration)
3. [Architecture](#architecture)
4. [Commandes CLI](#commandes-cli)
5. [Formats de données](#formats-de-données)
6. [Interface Git-like](#interface-git-like)
7. [Tests et qualité](#tests-et-qualité)
8. [Développement](#développement)
9. [Roadmap](#roadmap)

---

## 🎯 Vue d'ensemble

**Comptes** est un outil CLI minimal pour la gestion de finances personnelles, développé en Go. Il offre une interface Git-like avec un audit trail complet et des fonctionnalités avancées pour le suivi des transactions financières.

### ✨ Fonctionnalités principales

- **Gestion des transactions** : Ajout, édition, suppression avec audit trail
- **Interface Git-like** : Messages obligatoires, undo intelligent, historique complet
- **Formats multiples** : Text, CSV, JSON pour l'intégration avec d'autres outils
- **Architecture extensible** : Prête pour SQLite, REST API, etc.
- **Tests complets** : 28 tests automatiques avec edge cases

### 🎯 Objectifs

- **Simplicité** : Interface CLI intuitive et familière
- **Traçabilité** : Audit trail complet de toutes les modifications
- **Extensibilité** : Architecture en couches pour l'évolution future
- **Qualité** : Tests complets et validation automatique

---

## 🚀 Installation et configuration

### Prérequis

- **Go 1.21+** : Pour la compilation
- **Git** : Pour le contrôle de version
- **Terminal** : Interface en ligne de commande

### Installation

```bash
# Cloner le repository
git clone <repository-url>
cd comptes

# Compiler l'application
go build -o comptes cmd/comptes/main.go

# Initialiser le projet
./comptes init
```

### Configuration initiale

L'initialisation crée automatiquement :

```
comptes/
├── config/
│   └── config.yaml          # Configuration des comptes, catégories, tags
├── data/
│   ├── transactions.json     # Transactions financières
│   ├── accounts.json         # Comptes bancaires
│   ├── categories.json       # Catégories de transactions
│   └── tags.json            # Tags de transactions
└── comptes                  # Exécutable
```

### Configuration par défaut

Le fichier `config/config.yaml` contient :

```yaml
accounts:
  - id: "BANQUE"
    name: "Compte Courant Principal"
    type: "checking"
    currency: "EUR"
    initial_balance: 1500.00
    is_active: true

categories:
  - code: "ALM"
    name: "Alimentation"
    description: "Courses et repas"
  - code: "SLR"
    name: "Salaire"
    description: "Revenus professionnels"
  - code: "LGT"
    name: "Logement"
    description: "Loyer, charges, etc."

tags:
  - code: "URG"
    name: "Urgent"
    description: "Transaction urgente"
  - code: "REC"
    name: "Récurrent"
    description: "Transaction récurrente"
```

---

## 🏗️ Architecture

### Structure en couches

```
┌─────────────────┐
│   CLI Layer     │  ← Interface utilisateur
├─────────────────┤
│  Service Layer  │  ← Logique métier
├─────────────────┤
│  Storage Layer  │  ← Abstraction persistance
├─────────────────┤
│  Domain Layer   │  ← Modèles de données
└─────────────────┘
```

### Packages

- **`cmd/comptes/`** : Point d'entrée CLI
- **`internal/service/`** : Logique métier des transactions
- **`internal/storage/`** : Interface et implémentation JSON
- **`internal/domain/`** : Modèles de données
- **`internal/config/`** : Gestion de la configuration

### Modèles de données

#### Transaction
```go
type Transaction struct {
    ID          string    `json:"id"`
    Account     string    `json:"account"`
    Date        time.Time `json:"date"`
    Amount      float64   `json:"amount"`
    Description string    `json:"description"`
    Categories  []string  `json:"categories"`
    Tags        []string  `json:"tags"`
    IsActive    bool      `json:"is_active"`
    EditComment string    `json:"edit_comment,omitempty"`
    ParentID    string    `json:"parent_id,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

#### Account
```go
type Account struct {
    ID             string    `json:"id" yaml:"id"`
    Name           string    `json:"name" yaml:"name"`
    Type           string    `json:"type" yaml:"type"`
    Currency       string    `json:"currency" yaml:"currency"`
    InitialBalance float64   `json:"initial_balance" yaml:"initial_balance"`
    IsActive       bool      `json:"is_active" yaml:"is_active"`
    CreatedAt      time.Time `json:"created_at" yaml:"created_at"`
}
```

---

## 💻 Commandes CLI

### Commandes de base

#### Initialisation
```bash
./comptes init
```
Crée la structure de fichiers et la configuration par défaut.

#### Ajout de transactions
```bash
./comptes add '{"account": "BANQUE", "amount": -25.50, "description": "Courses", "categories": ["ALM"]}'
./comptes add '{"account": "BANQUE", "amount": 1500, "description": "Salaire", "categories": ["SLR"], "date": "today"}'
```

#### Liste des transactions
```bash
./comptes list                    # Transactions actives
./comptes list --history         # Toutes les transactions
./comptes list --format csv      # Format CSV
./comptes list --format json     # Format JSON
```

#### Édition de transactions
```bash
./comptes edit <id> '{"amount": -30.00}' -m "Correction montant"
./comptes edit <id> '{"description": "Nouvelle description"}' --message "Fix typo"
```

#### Suppression de transactions
```bash
./comptes delete <id> -m "Transaction erronée"
./comptes delete <id> --message "Suppression test"
```

#### Undo intelligent
```bash
./comptes undo <id>  # Détecte automatiquement le type d'opération
```

#### Calcul de solde
```bash
./comptes balance
```

#### Migration des IDs
```bash
./comptes migrate  # Convertit les anciens IDs vers UUID courts
```

### Options avancées

#### IDs partiels
```bash
./comptes edit fd66 '{"amount": -30.00}' -m "Test"  # Utilise les 4 premiers caractères
```

#### Dates flexibles
```bash
./comptes add '{"account": "BANQUE", "amount": -25.50, "description": "Test", "categories": ["ALM"], "date": "today"}'
./comptes add '{"account": "BANQUE", "amount": -25.50, "description": "Test", "categories": ["ALM"], "date": "yesterday"}'
./comptes add '{"account": "BANQUE", "amount": -25.50, "description": "Test", "categories": ["ALM"], "date": "2024-01-15"}'
```

---

## 📊 Formats de données

### Format JSON

#### Transaction complète
```json
{
  "id": "fd6647d8",
  "account": "BANQUE",
  "date": "2024-01-15",
  "amount": -25.50,
  "description": "Courses",
  "categories": ["ALM"],
  "tags": ["URG"],
  "is_active": true,
  "edit_comment": "",
  "parent_id": "",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

#### Transaction avec historique
```json
{
  "id": "fd6647d8",
  "date": "2024-01-15",
  "amount": -25.50,
  "description": "Courses",
  "categories": ["ALM"],
  "tags": [],
  "is_active": false,
  "edit_comment": "Correction montant"
}
```

### Format CSV

#### En-têtes
```csv
id,date,amount,description,categories,tags
```

#### Avec historique
```csv
id,date,amount,description,categories,tags,is_active,edit_comment
```

#### Exemple
```csv
fd6647d8,2024-01-15,-25.50,Courses,ALM,URG,true,
```

### Format Text

#### Transaction active
```
- [fd6647d8] ✅ 2024-01-15: -25.50 EUR - Courses (Categories: [ALM]), Tags: [URG]
```

#### Transaction supprimée
```
- [fd6647d8] ❌ 2024-01-15: -25.50 EUR - Courses (Categories: [ALM]) | Edit: Correction montant
```

---

## 🔄 Interface Git-like

### Philosophie

Comptes adopte la philosophie Git pour la gestion des transactions :

- **Messages obligatoires** : Chaque modification doit être documentée
- **Audit trail complet** : Historique de toutes les modifications
- **Undo intelligent** : Annulation des opérations avec détection automatique
- **Relations parent-enfant** : Traçabilité des éditions

### Workflow Git-like

#### Ajout
```bash
./comptes add '{"account": "BANQUE", "amount": -25.50, "description": "Courses", "categories": ["ALM"]}'
# → Transaction créée avec ID unique
```

#### Édition
```bash
./comptes edit abc123 '{"amount": -30.00}' -m "Correction montant"
# → Ancienne transaction désactivée avec commentaire
# → Nouvelle transaction créée avec parent_id = abc123
```

#### Suppression
```bash
./comptes delete def456 -m "Transaction erronée"
# → Transaction désactivée avec commentaire
```

#### Undo
```bash
./comptes undo abc123  # Undo edit : restaure parent, désactive enfant
./comptes undo def456  # Undo delete : restaure la transaction
./comptes undo ghi789  # Undo add : désactive la transaction
```

### Audit trail

Chaque modification est tracée :

1. **Transaction originale** : `is_active: true`, `edit_comment: ""`
2. **Après édition** : 
   - Ancienne : `is_active: false`, `edit_comment: "Correction montant"`
   - Nouvelle : `is_active: true`, `parent_id: "abc123"`
3. **Après undo** : Retour à l'état précédent

---

## 🧪 Tests et qualité

### Suite de tests

#### Tests automatiques
```bash
./test_mvp.sh  # 28 tests complets
```

#### Tests unitaires
```bash
go test ./...  # Tests des packages internes
```

#### Tests d'intégration
```bash
# Exécutés automatiquement dans pre-commit hooks
```

### Couverture de test

#### Tests de base (100% requis)
- ✅ Initialisation
- ✅ Ajout de transactions
- ✅ Liste avec tous les formats
- ✅ Édition avec messages
- ✅ Suppression avec messages
- ✅ Undo intelligent
- ✅ Calcul de solde

#### Tests d'erreurs (90% requis)
- ✅ Gestion des erreurs de validation
- ✅ Gestion des erreurs de commande
- ✅ Gestion des IDs inexistants/ambigus
- ✅ Gestion des dates invalides
- ✅ Gestion des fichiers corrompus

#### Tests de cohérence (95% requis)
- ✅ Audit trail complet
- ✅ Relations parent-enfant correctes
- ✅ Undo en chaîne fonctionne
- ✅ Migration des IDs fonctionne

#### Tests d'intégration (80% requis)
- ✅ Format CSV compatible Nushell
- ✅ Format JSON compatible jq
- ✅ Performance acceptable (100+ transactions)

### Pre-commit hooks

Validation automatique à chaque commit :

1. **`go mod tidy`** : Nettoyage des dépendances
2. **`go fmt`** : Formatage du code
3. **`go vet`** : Analyse statique
4. **`go test`** : Tests unitaires
5. **`go build`** : Compilation
6. **Tests d'intégration** : Validation du workflow complet

---

## 🔧 Développement

### Structure du projet

```
comptes/
├── cmd/
│   └── comptes/
│       └── main.go              # Point d'entrée CLI
├── internal/
│   ├── config/
│   │   ├── config.go           # Gestion configuration
│   │   └── config_test.go      # Tests configuration
│   ├── domain/
│   │   └── models.go           # Modèles de données
│   ├── service/
│   │   ├── transaction.go      # Logique métier
│   │   └── transaction_test.go # Tests service
│   └── storage/
│       ├── interface.go        # Interface storage
│       ├── json_storage.go     # Implémentation JSON
│       └── json_storage_test.go # Tests storage
├── docs/
│   ├── mvp.md                 # Spécifications MVP
│   ├── cli-commands.md        # Documentation CLI
│   ├── data-models.md         # Modèles de données
│   └── file-structure.md      # Structure fichiers
├── .git/
│   └── hooks/
│       └── pre-commit         # Hook de validation
├── test_mvp.sh               # Script de test
├── TEST_PLAN.md             # Plan de test complet
├── TODO.md                  # Progression du projet
├── README.md                # Documentation principale
├── SETUP.md                 # Guide d'installation
└── go.mod                   # Dépendances Go
```

### Commandes de développement

#### Compilation
```bash
go build -o comptes cmd/comptes/main.go
```

#### Tests
```bash
go test ./...                    # Tests unitaires
./test_mvp.sh                   # Tests d'intégration
```

#### Formatage
```bash
go fmt ./...                    # Formatage automatique
go vet ./...                    # Analyse statique
```

#### Dépendances
```bash
go mod tidy                     # Nettoyage
go mod download                 # Téléchargement
```

### Standards de code

- **Formatage** : `go fmt` automatique
- **Analyse** : `go vet` pour la qualité
- **Tests** : Couverture complète requise
- **Documentation** : Commentaires Go standards
- **Commits** : Messages descriptifs avec emojis

---

## 🚀 Roadmap

### MVP Complet ✅

- ✅ Architecture en couches
- ✅ CLI de base (init, add, list, edit, delete, undo, balance)
- ✅ Interface Git-like avec audit trail
- ✅ Formats multiples (text, csv, json)
- ✅ Tests complets (28 tests automatiques)
- ✅ Configuration par défaut
- ✅ Pre-commit hooks avec validation

### Fonctionnalités avancées (v2)

#### 🔄 Go routines & asynchrone
- **Analytics** : Calculs en parallèle sur plusieurs comptes
- **Import/Export** : Traitement de gros fichiers CSV/JSON
- **Validation** : Vérification des catégories/tags en parallèle
- **Cache** : Mise à jour asynchrone des soldes

#### 🎭 Mode transactionnel avec contexte
- **Contexte partagé** : `comptes account BANQUE` → `comptes category ALM`
- **Transaction atomique** : `comptes commit` (tout ou rien)
- **Moins verbeux** : Plus besoin de JSON pour chaque transaction
- **Validation groupée** : Vérification à la fin

#### 📅 Gestion des dates avancées
- **Formats étendus** : `last week`, `next month`, etc.
- **Calculs de dates** : `+7 days`, `-1 month`
- **Périodes** : Filtrage par mois, trimestre, année

#### 🧮 Calculs dans les requêtes
- **Expressions** : `{45.00 - 12.00}` pour les calculs
- **Références** : Montants dynamiques entre transactions
- **Validation** : Vérification des calculs et références

#### 💼 Fonctionnalités métier avancées
- **Support complet multi-comptes** avec transferts
- **Gestion des catégories** (CRUD via CLI)
- **Gestion des tags** (CRUD via CLI)
- **Règles de validation avancées** (catégories existantes, etc.)
- **Budgets et prévisions** : Gestion des budgets mensuels
- **Économies** : Suivi des objectifs d'épargne

### Évolutions techniques

#### 🗄️ Storage avancé
- **SQLite** : Base de données locale
- **REST API** : Intégration avec services externes
- **Synchronisation** : Multi-appareils
- **Backup** : Sauvegarde automatique

#### 🎨 Interface utilisateur
- **TUI** : Interface textuelle interactive
- **GUI** : Interface graphique simple
- **Web** : Interface web locale
- **Mobile** : Application mobile

#### 🔌 Intégrations
- **Banks** : Import automatique des relevés
- **APIs** : Intégration avec services financiers
- **Plugins** : Système de plugins
- **Export** : Formats d'export avancés

---

## 📞 Support et contribution

### Documentation

- **README.md** : Vue d'ensemble et installation
- **SETUP.md** : Guide de configuration
- **TEST_PLAN.md** : Plan de test complet
- **docs/** : Documentation détaillée

### Tests

- **`./test_mvp.sh`** : Tests complets
- **`go test ./...`** : Tests unitaires
- **Pre-commit hooks** : Validation automatique

### Développement

- **Architecture extensible** : Prête pour l'évolution
- **Tests complets** : Qualité assurée
- **Documentation** : Code bien documenté
- **Standards** : Respect des conventions Go

---

*Documentation mise à jour : 28 octobre 2025*
