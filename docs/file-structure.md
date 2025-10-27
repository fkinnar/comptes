# Structure du projet

## Organisation des fichiers

```
comptes/
├── README.md                    # Documentation principale
├── go.mod                       # Module Go
├── go.sum                       # Checksums des dépendances
├── cmd/
│   └── comptes/
│       └── main.go             # Point d'entrée de l'application
├── internal/
│   ├── domain/
│   │   ├── transaction.go      # Modèle Transaction
│   │   ├── comment.go           # Modèle Comment
│   │   ├── category.go          # Modèle Category
│   │   ├── tag.go               # Modèle Tag
│   │   ├── budget.go            # Modèle Budget
│   │   └── savings.go           # Modèle SavingsGoal
│   ├── storage/
│   │   ├── interface.go        # Interface Storage
│   │   ├── json/
│   │   │   ├── json_storage.go # Implémentation JSON
│   │   │   └── json_models.go  # Structures JSON
│   │   └── sqlite/             # Future implémentation SQLite
│   ├── service/
│   │   ├── transaction.go      # Service Transaction
│   │   ├── comment.go           # Service Comment
│   │   ├── budget.go            # Service Budget
│   │   ├── forecast.go          # Service Forecast
│   │   └── validation.go        # Validation des données
│   └── cli/
│       ├── commands.go         # Commandes CLI
│       ├── add.go              # Commande add
│       ├── list.go             # Commande list
│       ├── edit.go             # Commande edit
│       ├── delete.go           # Commande delete
│       ├── budget.go           # Commande budget
│       ├── forecast.go         # Commande forecast
│       ├── categories.go       # Commande categories
│       ├── tags.go             # Commande tags
│       ├── savings.go          # Commande savings
│       └── comments.go         # Commande comments
├── data/                        # Données par défaut
│   ├── transactions.json
│   ├── comments.json
│   ├── categories.json
│   ├── tags.json
│   ├── budgets.json
│   └── savings.json
├── config/                      # Configuration
│   └── config.yaml
└── docs/                        # Documentation
    ├── data-models.md
    ├── cli-commands.md
    └── file-structure.md
```

## Détail des packages

### `cmd/comptes/`
Point d'entrée de l'application. Contient uniquement `main.go` qui :
- Parse les arguments de ligne de commande
- Initialise les services
- Démarre l'exécution des commandes

### `internal/domain/`
Contient les modèles métier (structs Go) :
- `Transaction` : Représentation d'une transaction
- `Category` : Représentation d'une catégorie
- `Tag` : Représentation d'un tag
- `Budget` : Représentation d'un budget mensuel
- `SavingsGoal` : Représentation d'un objectif d'épargne

### `internal/storage/`
Abstraction du stockage des données :
- `interface.go` : Interface `Storage` avec méthodes CRUD
- `json/` : Implémentation JSON
- `sqlite/` : Future implémentation SQLite

### `internal/service/`
Logique métier de l'application :
- `transaction.go` : Gestion des transactions
- `budget.go` : Gestion du budget
- `forecast.go` : Calculs de prévisions
- `validation.go` : Validation des données

### `internal/cli/`
Interface en ligne de commande :
- `commands.go` : Configuration des commandes
- Fichiers séparés par commande pour la lisibilité

### `data/`
Fichiers de données JSON :
- `transactions.json` : Liste des transactions
- `categories.json` : Définitions des catégories
- `tags.json` : Définitions des tags
- `budgets.json` : Budgets mensuels
- `savings.json` : Objectifs d'épargne

### `config/`
Configuration YAML :
- `config.yaml` : Configuration générale

## Conventions de nommage

- **Fichiers** : `snake_case.go`
- **Packages** : `lowercase`
- **Types** : `PascalCase`
- **Fonctions** : `PascalCase` (exportées), `camelCase` (privées)
- **Variables** : `camelCase`

## Gestion des erreurs

- Utilisation des erreurs Go standard
- Messages d'erreur explicites
- Validation des données avant enregistrement
- Pas de panique, gestion gracieuse des erreurs

## Tests

Structure de tests à ajouter plus tard :
```
internal/
├── domain/
│   └── transaction_test.go
├── storage/
│   └── json/
│       └── json_storage_test.go
└── service/
    └── transaction_test.go
```
