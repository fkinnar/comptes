# Commandes CLI

## Commandes principales

### `comptes add`

Ajouter une nouvelle transaction

```bash
# Ajouter une dépense (paramètres individuels)
comptes add --amount -25.50 --desc "Achat gâteau" --categories alimentation,cadeau --tags anniversaire

# Ajouter une transaction (JSON direct)
comptes add --json '{
  "account_id": "uuid",
  "amount": -25.50,
  "description": "Achat gâteau",
  "categories": ["alimentation", "cadeau"],
  "tags": ["anniversaire"],
  "date": "2024-01-15T10:30:00Z"
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

# Ajouter depuis un fichier JSON
comptes add --file transactions.json

# Virement entre comptes (JSON)
comptes add --json '{
  "from_account_id": "uuid_source",
  "to_account_id": "uuid_dest", 
  "amount": 100.00,
  "description": "Virement vers épargne",
  "categories": ["virement"],
  "tags": ["epargne"]
}'
```

**Options :**

- `--json` / `-j` : JSON de la transaction ou tableau de transactions
- `--file` / `-f` : Fichier JSON contenant une transaction ou un tableau
- `--amount` / `-a` : Montant (compatible avec paramètres individuels)
- `--desc` / `-d` : Description (compatible avec paramètres individuels)
- `--categories` / `-c` : Catégories (compatible avec paramètres individuels)
- `--tags` / `-t` : Tags (compatible avec paramètres individuels)
- `--date` : Date (défaut: aujourd'hui)

**Note :** Le JSON accepte soit un objet (transaction unique) soit un tableau (plusieurs transactions). Facilite le scripting et l'import en masse.

### `comptes list`

Lister les transactions

```bash
# Lister les transactions actives
comptes list

# Voir l'historique complet (y compris les supprimées)
comptes list --history

# Filtrer par catégorie
comptes list --category alimentation

# Filtrer par période
comptes list --from 2024-01-01 --to 2024-01-31

# Filtrer par montant
comptes list --min-amount 50.00
```

**Options :**

- `--history` / `-h` : Afficher l'historique complet (transactions supprimées incluses)
- `--category` / `-c` : Filtrer par catégorie
- `--tag` / `-t` : Filtrer par tag
- `--from` / `-f` : Date de début
- `--to` : Date de fin
- `--min-amount` : Montant minimum
- `--max-amount` : Montant maximum

### `comptes edit`

Modifier une transaction existante

```bash
# Modifier avec paramètres individuels
comptes edit -id uuid -amount 30.00 -desc "Nouvelle description"

# Modifier avec JSON
comptes edit -id uuid -json '{
  "amount": 30.00,
  "description": "Nouvelle description",
  "categories": ["alimentation", "urgent"]
}'

# Modifier depuis un fichier
comptes edit -id uuid -file updated_transaction.json
```

**Options :**

- `-id` : ID de la transaction (requis)
- `-json` : JSON des modifications
- `-file` : Fichier JSON des modifications
- `-amount` : Nouveau montant (compatible avec paramètres individuels)
- `-desc` : Nouvelle description (compatible avec paramètres individuels)
- `-categories` : Nouvelles catégories (compatible avec paramètres individuels)
- `-tags` : Nouveaux tags (compatible avec paramètres individuels)
- `-date` : Nouvelle date (compatible avec paramètres individuels)

**Note :** L'édition désactive l'ancienne transaction et crée une nouvelle. L'historique complet est préservé.

### `comptes delete`

Supprimer une transaction (soft delete)

```bash
# Supprimer une transaction par ID
comptes delete -id uuid

# Supprimer avec confirmation
comptes delete -id uuid -confirm

# Supprimer plusieurs transactions
comptes delete -ids uuid1,uuid2,uuid3

# Supprimer avec commentaire
comptes delete -id uuid -comment "Erreur de saisie, transaction dupliquée"
```

**Options :**

- `-id` : ID de la transaction à supprimer
- `-ids` : IDs multiples séparés par des virgules
- `-confirm` : Confirmer la suppression
- `-force` : Supprimer sans confirmation
- `-comment` : Commentaire sur la suppression

**Note :** La suppression est un "soft delete" - la transaction est marquée comme `is_active: false` mais reste dans les données pour l'historique.

### `comptes budget`

Gérer le budget

```bash
# Voir le budget actuel
comptes budget

# Voir tous les budgets
comptes budget list

# Activer un budget
comptes budget activate -id uuid
```

**Sous-commandes :**

- `list` : Lister tous les budgets
- `activate` : Activer un budget spécifique

**Note :** Pour le MVP, les budgets sont configurés directement en JSON/YAML. Les commandes permettent seulement de les consulter et activer.

### `comptes forecast`

Calculer les prévisions

```bash
# Prévision pour les 30 prochains jours
comptes forecast -days 30

# Prévision jusqu'à la prochaine paie
comptes forecast -until-salary

# Prévision avec revenus moyens
comptes forecast -with-income
```

**Options :**

- `-days` : Nombre de jours à prévoir
- `-until-salary` : Jusqu'à la prochaine paie
- `-with-income` : Inclure les revenus moyens

### `comptes categories`

Gérer les catégories

```bash
# Lister toutes les catégories
comptes categories list

# Ajouter une catégorie
comptes categories add -code supermarche -name "Super marché" -parent alimentation

# Supprimer une catégorie
comptes categories remove supermarche
```

**Sous-commandes :**

- `list` : Lister les catégories
- `add` : Ajouter une catégorie
- `remove` : Supprimer une catégorie

**Note :** Une catégorie ne peut pas être supprimée si elle est utilisée dans des transactions. Les transactions avec des catégories "orphelines" restent valides.

### `comptes tags`

Gérer les tags

```bash
# Lister tous les tags
comptes tags list

# Ajouter un tag
comptes tags add -code anniversaire -name "Anniversaire" -parent evenement

# Supprimer un tag
comptes tags remove anniversaire
```

**Sous-commandes :**

- `list` : Lister les tags
- `add` : Ajouter un tag
- `remove` : Supprimer un tag

**Note :** Un tag ne peut pas être supprimé s'il est utilisé dans des transactions. Les transactions avec des tags "orphelins" restent valides.

### `comptes savings`

Gérer les objectifs d'épargne

```bash
# Lister les objectifs
comptes savings list

# Ajouter un objectif
comptes savings add -name "Citerne mazout" -target 2000.00 -date 2024-06-01

# Ajouter de l'argent à un objectif
comptes savings add-money -id uuid -amount 100.00
```

**Sous-commandes :**

- `list` : Lister les objectifs
- `add` : Ajouter un objectif
- `add-money` : Ajouter de l'argent à un objectif
- `remove` : Supprimer un objectif

### `comptes comments`

Gérer les commentaires de transactions

```bash
# Voir les commentaires d'une transaction
comptes comments -transaction-id uuid

# Ajouter un commentaire à une transaction
comptes comments add -transaction-id uuid -action modified -comment "Correction après vérification"

# Lister tous les commentaires récents
comptes comments list -days 30
```

**Sous-commandes :**

- `add` : Ajouter un commentaire
- `list` : Lister les commentaires
- `-transaction-id` : ID de la transaction
- `-action` : Type d'action (created, modified, deleted)
- `-comment` : Texte du commentaire
- `-days` : Nombre de jours pour lister les commentaires récents

## Commandes utilitaires

### `comptes init`

Initialiser le projet

```bash
# Créer la structure de fichiers
comptes init

# Initialiser avec des catégories par défaut
comptes init --with-defaults

# Initialiser dans un répertoire spécifique
comptes init --config-dir /path/to/config
```

### `comptes status`

Afficher le statut général

```bash
# Vue d'ensemble
comptes status

# Statut détaillé
comptes status -detailed
```

## Format de sortie

- **Tableau** : Par défaut pour les listes
- **JSON** : Avec l'option `-json`
- **CSV** : Avec l'option `-csv`
