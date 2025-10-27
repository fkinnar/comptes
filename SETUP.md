# Setup Guide - Comptes

Ce guide explique comment configurer le projet Comptes pour la première utilisation.

## 📋 Prérequis

- Go 1.19+ installé
- Accès en écriture au répertoire du projet

## 🚀 Installation rapide

### 1. Compiler le projet
```bash
go build -o comptes cmd/comptes/main.go
```

### 2. Initialiser le projet
```bash
./comptes init
```

Le projet est maintenant prêt à être utilisé !

## ⚙️ Configuration minimale

### Fichier de configuration : `config/config.yaml`

Le fichier `config/config.yaml` doit contenir au minimum :

#### **1. Comptes (obligatoire)**
```yaml
accounts:
  - id: "account1"                    # ID unique du compte
    name: "Compte Courant Principal"  # Nom affiché
    type: "checking"                  # Type de compte
    currency: "EUR"                   # Devise
    initial_balance: 1500.00          # Solde initial
```

#### **2. Catégories (obligatoire)**
```yaml
categories:
  - code: "alimentation"             # Code unique
    name: "Alimentation"              # Nom affiché
    description: "Dépenses alimentaires"
  - code: "transport"
    name: "Transport"
    description: "Dépenses de transport"
  - code: "logement"
    name: "Logement"
    description: "Dépenses de logement"
  - code: "salaire"
    name: "Salaire"
    description: "Revenus salariaux"
```

#### **3. Tags (optionnel)**
```yaml
tags:
  - code: "urgent"
    name: "Urgent"
    description: "Transaction urgente"
  - code: "recurrent"
    name: "Récurrent"
    description: "Transaction récurrente"
```

## 📝 Exemple complet

Voici un fichier `config/config.yaml` complet :

```yaml
# Configuration des comptes
accounts:
  - id: "account1"
    name: "Compte Courant Principal"
    type: "checking"
    currency: "EUR"
    initial_balance: 1500.00

# Configuration des catégories
categories:
  - code: "alimentation"
    name: "Alimentation"
    description: "Dépenses alimentaires"
  - code: "transport"
    name: "Transport"
    description: "Dépenses de transport"
  - code: "logement"
    name: "Logement"
    description: "Dépenses de logement"
  - code: "salaire"
    name: "Salaire"
    description: "Revenus salariaux"

# Configuration des tags
tags:
  - code: "urgent"
    name: "Urgent"
    description: "Transaction urgente"
  - code: "recurrent"
    name: "Récurrent"
    description: "Transaction récurrente"
```

## 🔧 Personnalisation

### Ajouter un nouveau compte
```yaml
accounts:
  - id: "account1"
    name: "Compte Courant Principal"
    type: "checking"
    currency: "EUR"
    initial_balance: 1500.00
  - id: "account2"                    # Nouveau compte
    name: "Compte Épargne"
    type: "savings"
    currency: "EUR"
    initial_balance: 5000.00
```

### Ajouter une nouvelle catégorie
```yaml
categories:
  - code: "alimentation"
    name: "Alimentation"
    description: "Dépenses alimentaires"
  - code: "loisirs"                   # Nouvelle catégorie
    name: "Loisirs"
    description: "Dépenses de loisirs"
```

### Ajouter un nouveau tag
```yaml
tags:
  - code: "urgent"
    name: "Urgent"
    description: "Transaction urgente"
  - code: "travail"                   # Nouveau tag
    name: "Travail"
    description: "Dépenses liées au travail"
```

## ⚠️ Règles importantes

### **Comptes**
- **`id`** : Obligatoire et unique
- **`name`** : Obligatoire
- **`type`** : Optionnel (défaut: "checking")
- **`currency`** : Optionnel (défaut: "EUR")
- **`initial_balance`** : Optionnel (défaut: 0.00)

### **Catégories**
- **`code`** : Obligatoire et unique
- **`name`** : Obligatoire
- **`description`** : Optionnel

### **Tags**
- **`code`** : Obligatoire et unique
- **`name`** : Obligatoire
- **`description`** : Optionnel

## 🚨 Erreurs courantes

### **1. Fichier config manquant**
```
Error: config file not found: /path/to/config/config.yaml
```
**Solution** : Créer le fichier `config/config.yaml` avec la structure minimale

### **2. ID de compte manquant**
```
Error: account ID is required
```
**Solution** : Ajouter un `id` unique à chaque compte

### **3. Code de catégorie manquant**
```
Error: category not found: alimentation
```
**Solution** : Vérifier que la catégorie existe dans la config

### **4. Code de tag manquant**
```
Error: tag not found: urgent
```
**Solution** : Vérifier que le tag existe dans la config

## 🔄 Après modification de la config

Après avoir modifié `config/config.yaml` :

1. **Supprimer les données existantes** (optionnel)
   ```bash
   rm -rf data/
   ```

2. **Réinitialiser le projet**
   ```bash
   ./comptes init
   ```

3. **Vérifier la configuration**
   ```bash
   ./comptes balance
   ```

## 📚 Prochaines étapes

Une fois la configuration terminée, vous pouvez :

1. **Ajouter des transactions**
   ```bash
   ./comptes add '{"account_id":"account1","amount":-25.50,"description":"Achat","categories":["alimentation"]}'
   ```

2. **Voir les transactions**
   ```bash
   ./comptes list
   ```

3. **Vérifier les soldes**
   ```bash
   ./comptes balance
   ```

## 🆘 Besoin d'aide ?

- Consultez le [README.md](README.md) pour l'utilisation générale
- Consultez [docs/cli-commands.md](docs/cli-commands.md) pour les commandes détaillées
- Consultez [docs/data-models.md](docs/data-models.md) pour la structure des données

