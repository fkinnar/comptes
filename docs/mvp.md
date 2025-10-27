# MVP - Version Minimale Viable

## Objectif
Créer une première version fonctionnelle avec les fonctionnalités essentielles, en laissant de côté les aspects complexes pour plus tard.

## Fonctionnalités MVP (Version 1.0)

### ✅ **Core - Transactions**
- Ajouter une transaction (`add`)
- Lister les transactions (`list`)
- Modifier une transaction (`edit`) - soft delete + création
- Supprimer une transaction (`delete`) - soft delete
- **Gestion des comptes multiples** (nouveau !)

### ✅ **Core - Données de base**
- Catégories simples (pas de hiérarchie au début)
- Tags simples (pas de hiérarchie au début)
- Stockage JSON basique

### ❌ **Reporté pour plus tard**
- Commentaires sur les transactions
- Budget mensuel
- Prévisions financières
- Objectifs d'épargne
- Validation avancée des données
- Création automatique de catégories/tags
- Hiérarchies de catégories/tags

## Gestion des comptes multiples

### Structure des comptes
```json
{
  "id": "uuid",
  "name": "Compte Courant Principal",
  "type": "checking",
  "currency": "EUR",
  "initial_balance": 1500.00,
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z"
}
```

### Transactions avec comptes
```json
{
  "id": "uuid",
  "account_id": "uuid",
  "date": "2024-01-15T10:30:00Z",
  "amount": -25.50,
  "description": "Achat gâteau",
  "categories": ["alimentation"],
  "tags": ["anniversaire"],
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Virements entre comptes (2 transactions liées)
```json
// Transaction de débit (compte source)
{
  "id": "uuid1",
  "account_id": "uuid_source",
  "date": "2024-01-15T10:30:00Z",
  "amount": -100.00,
  "description": "Virement vers épargne",
  "categories": ["virement"],
  "tags": ["epargne"],
  "transfer_id": "uuid_transfer",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}

// Transaction de crédit (compte destination)
{
  "id": "uuid2",
  "account_id": "uuid_dest",
  "date": "2024-01-15T10:30:00Z",
  "amount": 100.00,
  "description": "Virement depuis compte courant",
  "categories": ["virement"],
  "tags": ["epargne"],
  "transfer_id": "uuid_transfer",
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

## Commandes MVP

### `comptes accounts`
```bash
# Lister les comptes avec soldes
comptes accounts list

# Ajouter un compte
comptes accounts add -name "Compte Courant" -type checking -initial-balance 1500.00

# Supprimer un compte
comptes accounts remove -id uuid
```

### `comptes add` (modifié)
```bash
# Ajouter une transaction (paramètres individuels)
comptes add --account-id uuid --amount -25.50 --desc "Achat" --categories alimentation

# Ajouter une transaction (JSON direct)
comptes add --json '{
  "account_id": "uuid",
  "amount": -25.50,
  "description": "Achat gâteau",
  "categories": ["alimentation"]
}'

# Ajouter plusieurs transactions en une fois
comptes add --json '[
  {
    "account_id": "uuid",
    "amount": -25.50,
    "description": "Achat gâteau",
    "categories": ["alimentation"]
  },
  {
    "account_id": "uuid",
    "amount": 2500.00,
    "description": "Salaire",
    "categories": ["salaire"]
  }
]'

# Virement entre comptes (JSON)
comptes add --json '{
  "from_account_id": "uuid_source",
  "to_account_id": "uuid_dest",
  "amount": 100.00,
  "description": "Virement"
}'
```

### `comptes list` (simplifié)
```bash
# Lister toutes les transactions actives
comptes list

# Voir l'historique complet
comptes list --history

# Voir les soldes des comptes
comptes list --balances
```

**Note :** Les filtres avancés (par catégorie, montant, période) sont laissés à des outils externes comme Nushell pour le MVP. L'outil se contente de sortir du JSON propre.

## Simplifications MVP

### Validation minimale
- Vérifier que les comptes existent
- Vérifier que les catégories/tags existent (pas de création auto)
- Montants positifs
- Dates cohérentes

### Pas de fonctionnalités avancées
- Pas de commentaires
- Pas de budget
- Pas de prévisions
- Pas de hiérarchies
- Pas de création automatique
- **Pas de filtres avancés** (délégués à Nushell/outils externes)

### Structure de fichiers simplifiée
```
comptes/                    # Répertoire de l'exécutable
├── comptes                 # Exécutable
├── data/                   # Données (à côté de l'exécutable)
│   ├── accounts.json
│   ├── transactions.json
│   ├── categories.json
│   ├── tags.json
│   └── budgets.json
└── config/
    └── config.yaml
```

### Configuration des budgets
Les budgets sont configurés directement en JSON/YAML avec des dates de validité :
- `valid_from` : Date de début
- `valid_to` : Date de fin (optionnel)
- Un seul budget actif à la fois

### Recherche des fichiers
- **MVP** : Fichiers recherchés à côté de l'exécutable
- **Évolution** : Migration vers `~/.config/comptes/` plus tard
- **Flexibilité** : Possibilité de spécifier un répertoire avec `--config-dir`

## Roadmap post-MVP

### Version 1.1
- Commentaires sur les transactions
- Validation avancée
- **Filtres avancés** (nice to have)
- **Migration vers `~/.config/comptes/`**

### Version 1.2
- Budget mensuel
- Hiérarchies de catégories/tags

### Version 1.3
- Prévisions financières
- Objectifs d'épargne

### Version 2.0
- Migration vers SQLite
- Interface web (optionnel)

## TODO - Fonctionnalités avancées des comptes

### Types de comptes configurables
- Définir des types personnalisés (comme les catégories)
- Règles par type de compte :
  - Limitation à certaines catégories (ex: ticket resto → alimentation uniquement)
  - Source de virement possible (ex: pas de virement depuis carte de crédit)
  - Solde négatif autorisé et limites
  - Frais de découvert
  - Intérêts créditeurs/débiteurs

### Gestion multi-devises
- Support de plusieurs devises
- Taux de change
- Conversion automatique
- Comptes en devises étrangères

### Règles métier avancées
- Validation des virements selon les types de comptes
- Alertes de solde faible
- Limites de transaction par type de compte
- Rapprochement bancaire automatique
