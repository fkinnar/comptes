# Modèles de données

## Transactions

### Structure JSON
```json
{
  "id": "uuid",
  "date": "2024-01-15T10:30:00Z",
  "amount": -25.50,
  "description": "Achat gâteau anniversaire",
  "categories": ["alimentation", "cadeau"],
  "tags": ["anniversaire", "urgent"],
  "is_active": true,
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

### Signification des montants
- **Montant positif** : Revenus/entrées d'argent
- **Montant négatif** : Dépenses/sorties d'argent

### Champs
- `is_active` : `true` pour les transactions actives, `false` pour les supprimées
- `created_at` : Date de création
- `updated_at` : Date de dernière modification

## Commentaires de transactions

### Structure JSON
```json
{
  "id": "uuid",
  "transaction_id": "uuid",
  "action": "created|modified|deleted",
  "comment": "Commentaire sur l'action",
  "created_at": "2024-01-15T10:30:00Z",
  "created_by": "user" // Optionnel pour l'avenir
}
```

### Types d'actions
- `created` : Transaction créée
- `modified` : Transaction modifiée
- `deleted` : Transaction supprimée

### Règles
- Un commentaire est lié à une transaction via `transaction_id`
- Plusieurs commentaires peuvent exister pour une même transaction
- Les commentaires ne sont jamais supprimés (historique complet)
- Chaque action (création, modification, suppression) peut avoir un commentaire

## Catégories

### Structure JSON
```json
{
  "code": "alimentation",
  "name": "Alimentation",
  "parent": null,
  "children": ["supermarche", "restaurant"],
  "description": "Toutes les dépenses liées à l'alimentation"
}
```

### Hiérarchie
- Catégories parentes (ex: "alimentation")
- Catégories enfants (ex: "supermarche", "restaurant")
- Support des multi-niveaux

### Règles
- Le `code` est unique et référencé dans les transactions
- Si une catégorie est utilisée dans des transactions, elle ne peut pas être supprimée
- Les transactions ne peuvent pas être créées avec des catégories inexistantes (sauf avec `-create-missing`)
- Si le code n'existe plus, la transaction reste mais avec une catégorie "orpheline"

## Tags

### Structure JSON
```json
{
  "code": "anniversaire",
  "name": "Anniversaire",
  "parent": "evenement",
  "children": [],
  "description": "Événements liés aux anniversaires"
}
```

### Hiérarchie
- Tags parentes (ex: "evenement")
- Tags enfants (ex: "anniversaire", "noel")
- Support des multi-niveaux

### Règles
- Le `code` est unique et référencé dans les transactions
- Si un tag est utilisé dans des transactions, il ne peut pas être supprimé
- Les transactions ne peuvent pas être créées avec des tags inexistants (sauf avec `-create-missing`)
- Si le code n'existe plus, la transaction reste mais avec un tag "orphelin"

## Budget

### Structure JSON
```json
{
  "id": "uuid",
  "name": "Budget 2024",
  "version": "1.0",
  "valid_from": "2024-01-01",
  "valid_to": "2024-12-31",
  "fixed_expenses": [
    {
      "category": "logement",
      "amount": 800.00,
      "description": "Loyer"
    }
  ],
  "envelopes": [
    {
      "category": "alimentation",
      "amount": 300.00,
      "description": "Budget courses"
    }
  ],
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z"
}
```

### Champs
- `valid_from` : Date de début de validité du budget
- `valid_to` : Date de fin de validité du budget (optionnel, null = sans limite)
- `is_active` : Budget actuellement utilisé

## Objectifs d'épargne

### Structure JSON
```json
{
  "id": "uuid",
  "name": "Citerne mazout",
  "target_amount": 2000.00,
  "current_amount": 500.00,
  "target_date": "2024-06-01",
  "category": "maison",
  "description": "Économies pour remplir la citerne"
}
```

## Configuration (YAML)

### Structure
```yaml
# config.yaml
categories:
  - code: alimentation
    name: Alimentation
    children:
      - supermarche
      - restaurant
  
tags:
  - code: evenement
    name: Événement
    children:
      - anniversaire
      - noel

settings:
  currency: "EUR"
  date_format: "2006-01-02"
  default_categories:
    - alimentation
    - transport
```
