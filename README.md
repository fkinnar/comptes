# Comptes - Gestionnaire de comptes personnel

Un outil en ligne de commande minimal pour gérer ses comptes personnels, écrit en Go.

## 🚀 Démarrage rapide

### Compilation
```bash
# Compiler l'exécutable
go build -o comptes cmd/comptes/main.go

# Ou utiliser go run pour tester
go run cmd/comptes/main.go init
```

### Configuration et initialisation
```bash
# 1. Configurer le projet (voir SETUP.md pour les détails)
# 2. Initialiser le projet (crée les fichiers de données)
./comptes init
```

### Utilisation
```bash
# Ajouter une transaction
./comptes add '{"account_id":"account1","amount":-25.50,"description":"Achat gâteau","categories":["alimentation"]}'

# Voir les transactions
./comptes list

# Voir les soldes
./comptes balance
```

> 📖 **Configuration requise** : Consultez [SETUP.md](SETUP.md) pour configurer les comptes, catégories et tags nécessaires.

## 📁 Structure du projet

```
comptes/
├── comptes                 # Exécutable compilé
├── go.mod                  # Module Go
├── cmd/
│   └── comptes/
│       └── main.go         # Point d'entrée de l'application
├── internal/               # Code interne (non importable)
│   ├── domain/
│   │   └── models.go       # Modèles de données (Account, Transaction, etc.)
│   ├── storage/
│   │   ├── interface.go     # Interface Storage
│   │   └── json_storage.go  # Implémentation JSON
│   └── service/
│       └── transaction.go   # Logique métier des transactions
├── data/                   # Données JSON (créé à l'exécution)
│   ├── accounts.json
│   ├── transactions.json
│   ├── categories.json
│   └── tags.json
├── config/
│   └── config.yaml         # Configuration YAML
└── docs/                   # Documentation
    ├── README.md
    ├── mvp.md
    ├── cli-commands.md
    ├── data-models.md
    └── file-structure.md
```

## 🏗️ Architecture

Le projet suit une architecture en couches pour permettre l'évolution du stockage des données :

### **cmd/comptes/** - Point d'entrée
- `main.go` : Point d'entrée de l'application
- Parse les arguments de ligne de commande
- Initialise les services et démarre l'exécution

### **internal/domain/** - Entités métier
- `models.go` : Définit les structures de données
- `Account` : Comptes bancaires
- `Transaction` : Transactions financières
- `Category` : Catégories de transactions
- `Tag` : Tags de transactions

### **internal/storage/** - Persistance des données
- `interface.go` : Interface `Storage` avec méthodes CRUD
- `json_storage.go` : Implémentation JSON (fichiers locaux)
- Extensible vers SQLite ou autres systèmes

### **internal/service/** - Logique métier
- `transaction.go` : Service pour gérer les transactions
- Validation des données
- Calcul des soldes
- Logique d'ajout/modification

### **data/** - Données persistantes
- Fichiers JSON créés automatiquement
- Un fichier par type de données
- Sauvegarde automatique à chaque modification

## 🔧 Développement

### Compilation
```bash
# Compiler
go build -o comptes cmd/comptes/main.go

# Compiler avec optimisations
go build -ldflags="-s -w" -o comptes cmd/comptes/main.go

# Cross-compilation (exemple pour Linux depuis macOS)
GOOS=linux GOARCH=amd64 go build -o comptes-linux cmd/comptes/main.go
```

### Tests
```bash
# Lancer les tests (quand ils existeront)
go test ./...

# Tests avec couverture
go test -cover ./...
```

### Modules Go
```bash
# Ajouter une dépendance
go get github.com/example/package

# Nettoyer les dépendances
go mod tidy

# Vérifier les dépendances
go mod verify
```

## 📊 Données

### Structure des fichiers JSON
- **accounts.json** : Comptes avec soldes initiaux
- **transactions.json** : Toutes les transactions
- **categories.json** : Catégories disponibles
- **tags.json** : Tags disponibles

### Configuration YAML
- **config/config.yaml** : Configuration par défaut
- Comptes, catégories et tags de base

## 🎯 Objectifs

- Enregistrer entrées et dépenses
- Créer et gérer un budget mensuel
- Gérer les dépenses récurrentes
- Planifier des économies pour des dépenses futures
- Calculer des prévisions financières

## 🏛️ Principes de design

- **Simplicité** : Utilisation des bibliothèques standard Go autant que possible
- **Extensibilité** : Architecture modulaire pour faciliter l'évolution
- **Lisibilité** : Configuration en YAML, données en JSON
- **Validation** : Pas d'enregistrement si les données ne sont pas valides

## 💾 Stockage des données

- **Format** : JSON pour les données, YAML pour la configuration
- **Évolutivité** : Interface abstraite pour permettre le passage à SQLite ou autres
- **Structure** : Fichiers séparés par type de données
- **Localisation** : Fichiers à côté de l'exécutable (MVP), migration vers `~/.config/comptes/` prévue

## 🎮 Commandes principales

- `init` : Initialiser le projet
- `add` : Ajouter une transaction
- `list` : Lister les transactions
- `balance` : Afficher les soldes des comptes
- `edit` : Modifier une transaction existante (à venir)
- `delete` : Supprimer une transaction (à venir)
- `budget` : Gérer le budget mensuel (à venir)
- `forecast` : Calculer les prévisions (à venir)
- `categories` : Gérer les catégories (à venir)
- `tags` : Gérer les mots-clés (à venir)
- `savings` : Gérer les objectifs d'épargne (à venir)
- `comments` : Gérer les commentaires de transactions (à venir)
