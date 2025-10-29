# Commandes CLI - Documentation complète

## 📋 Table des matières

1. [Commandes implémentées](#commandes-implémentées)
2. [Fonctionnalités à venir](#fonctionnalités-à-venir)

---

## ✅ Commandes implémentées

### `comptes init`

Initialise le projet et crée les fichiers de configuration et de données nécessaires.

```bash
# Initialiser le projet
comptes init
```

**Ce qui est créé :**
- `config/config.yaml` : Configuration par défaut (comptes, catégories, tags)
- `data/accounts.json` : Comptes bancaires
- `data/movements.json` : Mouvements financiers (vide)
- `data/categories.json` : Catégories de transactions
- `data/tags.json` : Tags de transactions

**Configuration par défaut :**
- 2 comptes : BANQUE (Compte Courant) et LIVRET (Compte Épargne)
- 10 catégories : Alimentation, Salaire, Logement, Transport, Santé, Loisirs, Vêtements, Télécom, Assurance, Éducation
- 4 tags : Urgent, Récurrent, Important, Professionnel

---

### `comptes add`

Ajoute un mouvement financier (transaction individuelle ou dans un batch).

Vous pouvez utiliser soit le format JSON, soit des flags pour ajouter un mouvement.

#### Format JSON

```bash
# Ajouter directement un mouvement
comptes add '{"account":"BANQUE","amount":-25.50,"description":"Courses","categories":["ALM"]}'

# Ajouter avec date relative
comptes add '{"account":"BANQUE","amount":-25.50,"description":"Courses","date":"today"}'

# Ajouter dans une transaction batch courante
comptes begin "Dépenses du mois"
comptes add '{"account":"BANQUE","amount":-25.50,"description":"Courses"}'
comptes add '{"account":"BANQUE","amount":-10.00,"description":"Pain"}'
comptes commit

# Ajouter dans une batch spécifique
comptes add '{"account":"BANQUE","amount":-25.50,"description":"Courses"}' <batch-id>

# Forcer l'ajout direct même si une batch est en cours
comptes add '{"account":"BANQUE","amount":-25.50,"description":"Courses"}' --immediate
comptes add -a BANQUE -m -25.50 -d "Courses" -i  # Version courte
```

#### Format Flags (plus simple pour l'usage quotidien)

```bash
# Ajouter directement un mouvement avec flags
comptes add -a BANQUE -m -25.50 -d "Courses" -c ALM
comptes add --account BANQUE --amount -25.50 --description "Courses" --categories "ALM,SLR" --tags "REC,URG"

# Ajouter avec date
comptes add -a BANQUE -m -25.50 -d "Courses" -o today
comptes add -a BANQUE -m -25.50 -d "Courses" --on yesterday

# Ajouter dans une batch courante
comptes begin "Dépenses du mois"
comptes add -a BANQUE -m -25.50 -d "Courses" -c ALM
comptes add -a BANQUE -m -10.00 -d "Pain" -c ALM
comptes commit

# Ajouter dans une batch spécifique
comptes add -a BANQUE -m -25.50 -d "Courses" <batch-id>

# Forcer l'ajout direct avec flags
comptes add -a BANQUE -m -25.50 -d "Courses" --immediate
comptes add -a BANQUE -m -25.50 -d "Courses" -i  # Version courte
```

**Format JSON :**
```json
{
  "account": "BANQUE",           // ID du compte (requis)
  "amount": -25.50,              // Montant (requis, négatif = dépense, positif = revenu)
  "description": "Courses",      // Description (requis)
  "categories": ["ALM"],         // Catégories (optionnel, codes)
  "tags": ["REC"],               // Tags (optionnel, codes)
  "date": "today"                // Date (optionnel, formats: today, yesterday, 2024-01-15)
}
```

**Flags disponibles :**
- `-a, --account <id>` : ID du compte (requis)
- `-m, --amount <value>` : Montant (requis, négatif = dépense, positif = revenu)
- `-d, --desc, --description <text>` : Description (requis)
- `-c, --categories <codes>` : Catégories (optionnel, séparées par virgules, ex: "ALM,SLR")
- `-t, --tags <codes>` : Tags (optionnel, séparés par virgules, ex: "REC,URG")
- `-o, --on, --date <date>` : Date (optionnel, formats: today, yesterday, 2024-01-15)
- `-i, --immediate` : Force l'ajout immédiat dans `movements.json`, même si une batch est en cours

**Détection automatique :**
- Si le premier argument commence par `{` ou `[`, le mode JSON est utilisé
- Sinon, le mode flags est utilisé

**Options communes :**
- Si `batch-id` est fourni (ou batch courante définie), le mouvement est ajouté au batch au lieu d'être directement enregistré
- Format de date flexible : `today`, `yesterday`, `tomorrow`, `2024-01-15`, `15/01/2024`

**Cas d'usage du flag `--immediate` :**
- Vous avez commencé une batch pour des dépenses mensuelles
- Vous avez besoin d'ajouter une transaction urgente qui doit être immédiatement visible
- Vous voulez ajouter quelques transactions en dehors de la batch sans avoir à rollback/commit la batch courante

---

### `comptes list`

Liste les données selon le type demandé.

```bash
# Lister les mouvements (défaut)
comptes list

# Lister avec historique complet (supprimés inclus)
comptes list --history

# Lister les catégories
comptes list --categories
comptes list --categories --format csv
comptes list --categories --format json

# Lister les tags
comptes list --tags
comptes list --tags --format csv

# Lister les comptes avec soldes
comptes list --accounts
comptes list --accounts --format csv

# Formats de sortie
comptes list --format text   # Format texte (défaut)
comptes list --format csv    # CSV compatible Nushell
comptes list --format json   # JSON pour scripting

# Afficher les codes au lieu des noms
comptes list --codes
```

**Options :**
- `--transactions` : Liste les mouvements (défaut)
- `--categories, -c` : Liste les catégories disponibles
- `--tags, -t` : Liste les tags disponibles
- `--accounts, -a` : Liste les comptes avec leurs soldes actuels
- `--history, -h` : Affiche tous les mouvements (y compris supprimés/édités)
- `--format <fmt>` : Format de sortie (`text`, `csv`, `json`)
- `--codes` : Affiche les codes au lieu des noms complets

**Informations affichées pour les comptes :**
- Nom du compte
- ID du compte
- Type de compte
- Devise
- Solde actuel (calculé à partir des mouvements)
- Solde initial (affiché si différent)

---

### `comptes edit`

Modifie un mouvement existant (soft delete + création d'un nouveau).

```bash
# Modifier un mouvement
comptes edit <id> '{"amount": -30.00}' -m "Correction montant"
comptes edit <id> '{"description": "Nouvelle description"}' --message "Fix typo"
comptes edit fd66 '{"amount": -30.00}' -m "Correction"  # ID partiel
```

**Options :**
- `-m, --message` : Message obligatoire expliquant la modification
- Support des IDs partiels (ex: `fd66` au lieu de l'UUID complet)

**Comportement :**
- L'ancien mouvement est marqué comme `is_active: false`
- Un nouveau mouvement est créé avec les modifications
- Relation parent-enfant préservée pour l'historique
- Audit trail complet

---

### `comptes delete`

Supprime un mouvement (soft delete).

```bash
# Suppression simple
comptes delete <id> -m "Transaction erronée"

# Suppression définitive
comptes delete <id> --hard -m "Dupliquée"
comptes delete <id> --hard --force -m "Dupliquée"  # Sans confirmation
```

**Options :**
- `-m, --message` : Message obligatoire expliquant la suppression
- `--hard` : Suppression définitive (pas de récupération possible)
- `-f, --force` : Bypasse la confirmation pour les opérations destructives
- Support des IDs partiels

---

### `comptes undo`

Annule la dernière opération sur un mouvement (add/edit/delete).

```bash
# Annuler l'opération (détection automatique du type)
comptes undo <id>

# Annulation définitive
comptes undo <id> --hard
comptes undo <id> --hard --force
```

**Options :**
- `--hard` : Suppression définitive au lieu de soft delete
- `-f, --force` : Bypasse la confirmation
- Support des IDs partiels

**Détection automatique :**
- Si mouvement créé par `edit` → annule l'édition (restaure l'ancien)
- Si mouvement supprimé → restaure le mouvement
- Si mouvement ajouté → supprime le mouvement

---

### `comptes balance`

Affiche les soldes de tous les comptes.

```bash
# Afficher les soldes
comptes balance
```

**Affichage :**
- Nom du compte
- Solde actuel (calculé à partir du solde initial + tous les mouvements actifs)
- Devise

---

### `comptes begin`

Commence une nouvelle transaction batch.

```bash
# Créer une batch sans description
comptes begin

# Créer une batch avec description
comptes begin "Dépenses du mois"
comptes begin "Transactions du 15 octobre"
```

**Comportement :**
- Crée une nouvelle batch avec un UUID unique
- La batch devient automatiquement la "batch courante"
- Les mouvements ajoutés avec `comptes add` (sans batch-id) utiliseront cette batch
- La batch est stockée dans `data/pending_transactions.json`

**Stockage de la batch courante :**
- Un fichier `.current_batch` dans `data/` stocke l'ID de la batch courante
- Permet d'utiliser `comptes commit` sans spécifier l'ID

---

### `comptes commit`

Commite une transaction batch (ajoute tous les mouvements dans `movements.json`).

```bash
# Commiter la batch courante
comptes commit

# Commiter une batch spécifique (ID partiel ou complet)
comptes commit abc12345
comptes commit abc12345-b623-4e51-933d-5f91ddac6136
```

**Options :**
- Si aucun `batch-id` n'est fourni, utilise la batch courante
- Support des IDs partiels (ex: `abc123` au lieu de l'UUID complet)
- Valide tous les mouvements avant le commit
- Si une validation échoue, le commit échoue entièrement (atomique)

**Comportement :**
- Tous les mouvements sont ajoutés dans `movements.json`
- La batch est déplacée dans `data/committed_transactions.json`
- La batch courante est effacée si c'était celle qui était commitée

---

### `comptes rollback`

Annule une transaction batch (supprime tous les mouvements de la batch).

```bash
# Rollback de la batch courante
comptes rollback

# Rollback d'une batch spécifique
comptes rollback abc12345
```

**Options :**
- Si aucun `batch-id` n'est fourni, utilise la batch courante
- Support des IDs partiels

**Comportement :**
- La batch est déplacée dans `data/rolled_back_transactions.json`
- Les mouvements ne sont jamais ajoutés dans `movements.json`
- La batch courante est effacée si c'était celle qui était rollbackée

---

## 🔄 Mode transactionnel

Le mode transactionnel permet de grouper plusieurs mouvements et de les valider ensemble.

### Workflow typique

```bash
# 1. Créer une batch
comptes begin "Dépenses du mois"

# 2. Ajouter plusieurs mouvements
comptes add '{"account":"BANQUE","amount":-45.50,"description":"Courses"}'
comptes add '{"account":"BANQUE","amount":-12.30,"description":"Pain"}'
comptes add '{"account":"BANQUE","amount":-80.00,"description":"Essence"}'

# 3. Commiter tous les mouvements d'un coup
comptes commit
```

### Avantages

- **Validation atomique** : Tous les mouvements sont validés ensemble
- **Rollback facile** : Si une erreur survient, `comptes rollback` annule tout
- **Batch courante** : Pas besoin de retenir l'UUID (utilise `.current_batch`)
- **IDs partiels** : Utilisez les premiers caractères de l'UUID (ex: `abc123`)

### Fichiers de stockage

- `data/pending_transactions.json` : Batches en attente
- `data/committed_transactions.json` : Historique des batches commitées
- `data/rolled_back_transactions.json` : Historique des batches rollbackées
- `data/.current_batch` : ID de la batch courante (fichier caché)

---

## ⏳ Fonctionnalités à venir

### Améliorations du mode transactionnel

#### Contexte partagé
Permettre de définir une fois les paramètres communs (compte, catégorie, tags) pour tous les mouvements d'une batch.

```bash
comptes begin "Dépenses du mois"
comptes account BANQUE           # Définit le compte pour cette batch
comptes category ALM             # Définit la catégorie par défaut
comptes add '{"amount": -25.50, "description": "Courses"}'  # Utilise le contexte
comptes add '{"amount": -10.00, "description": "Pain"}'     # Réutilise le contexte
comptes commit
```

**Statut :** Non implémenté

---

### Gestion des comptes

#### Commandes de gestion
```bash
comptes accounts add --id LIVRET --name "Compte Épargne" --type savings
comptes accounts edit LIVRET --name "Épargne long terme"
comptes accounts list
comptes accounts delete LIVRET
```

**Statut :** Non implémenté (les comptes sont actuellement uniquement dans la config)

---

### Gestion des catégories et tags

#### Commandes de gestion
```bash
comptes categories add --code VET --name "Vêtements"
comptes categories edit VET --name "Habillement"
comptes categories list
comptes categories delete VET

comptes tags add --code IMP --name "Important"
comptes tags list
comptes tags delete IMP
```

**Statut :** Non implémenté (les catégories/tags sont actuellement uniquement dans la config)
**Note :** La consultation existe déjà avec `comptes list --categories` et `comptes list --tags`

---

### Import/Export

#### Import CSV
```bash
comptes import --file bank_statement.csv
comptes import --csv "date,amount,description,category"
```

**Statut :** Non implémenté

#### Export amélioré
```bash
comptes export --format csv --from 2024-01-01 --to 2024-01-31
comptes export --format json --account BANQUE
```

**Statut :** Partiellement implémenté (via `comptes list --format csv/json`)

---

### Rapports et analytics

#### Rapports simples
```bash
comptes report --month 2024-01
comptes report --category ALM --from 2024-01-01
comptes report --account BANQUE --year 2024
```

**Statut :** Non implémenté

#### Tendances
```bash
comptes balance --trend
comptes balance --history 6  # 6 derniers mois
```

**Statut :** Non implémenté

---

### Multi-comptes avec transferts

#### Virements entre comptes
```bash
comptes add '{
  "account": "BANQUE",
  "amount": -100,
  "transfer_to": "LIVRET",
  "description": "Virement épargne"
}'
```

**Statut :** Non implémenté

---

### Budgets et prévisions

#### Gestion des budgets
```bash
comptes budget set --month 2024-01 --category ALM --amount 500
comptes budget show --month 2024-01
comptes budget compare --month 2024-01  # Comparer avec réel
```

**Statut :** Non implémenté

#### Prévisions
```bash
comptes forecast --days 30
comptes forecast --until-salary
comptes forecast --with-income
```

**Statut :** Non implémenté

---

### Objectifs d'épargne

#### Gestion des objectifs
```bash
comptes savings add --name "Vacances" --target 2000.00 --date 2024-06-01
comptes savings list
comptes savings add-money --id <goal-id> --amount 100.00
comptes savings progress
```

**Statut :** Non implémenté

---

### Filtres avancés pour `list`

#### Filtres par date
```bash
comptes list --from 2024-01-01 --to 2024-01-31
comptes list --month 2024-01
comptes list --year 2024
```

**Statut :** Non implémenté

#### Filtres par montant
```bash
comptes list --min-amount 50.00
comptes list --max-amount 500.00
comptes list --amount-range 50.00 500.00
```

**Statut :** Non implémenté

#### Filtres par compte
```bash
comptes list --account BANQUE
comptes list --accounts BANQUE,LIVRET
```

**Statut :** Non implémenté

---

### Recherche et navigation

#### Recherche dans les descriptions
```bash
comptes search "supermarket"
comptes search --account BANQUE "courses"
```

**Statut :** Non implémenté

---

### Installation système

#### Installation dans le PATH
```bash
make install  # Installe dans /usr/local/bin
```

**Statut :** Makefile existe mais nécessite configuration des chemins système

#### Respect XDG Base Directory
```bash
# ~/.config/comptes/config.yaml
# ~/.local/share/comptes/movements.json
```

**Statut :** Non implémenté (actuellement fichiers à côté de l'exécutable)

---

### Performance et optimisation

#### Snapshots de solde
Pour éviter de recalculer le solde depuis le début à chaque fois.

**Statut :** Non implémenté

#### Cache
Couche de cache entre Service et Storage pour optimiser les lectures.

**Statut :** Non implémenté

#### Indexation
Pour les recherches rapides par date, catégorie, etc.

**Statut :** Non implémenté

---

### Migration vers SQLite

Quand la structure sera validée à l'usage, migration vers SQLite pour de meilleures performances.

**Statut :** Non implémenté (architecture prête pour l'évolution)

---

## 📝 Notes de développement

### Architecture actuelle

- **Storage** : JSON files (extensible vers SQLite)
- **Service layer** : Validation, logique métier
- **CLI layer** : Interface utilisateur
- **Format des données** : JSON pour movements, batches (pending/committed/rolled_back)

### Structure des fichiers

```
data/
├── accounts.json              # Comptes bancaires
├── movements.json            # Mouvements financiers (anciennement transactions.json)
├── categories.json           # Catégories
├── tags.json                # Tags
├── pending_transactions.json # Batches en attente
├── committed_transactions.json # Batches commitées
├── rolled_back_transactions.json # Batches rollbackées
└── .current_batch            # Batch courante (caché)
```

### Modèles de données

- **Account** : Comptes avec solde initial
- **Transaction** : Mouvements individuels (dans movements.json)
- **TransactionBatch** : Groupes de mouvements (dans pending/committed/rolled_back)
- **Category** : Catégories de transactions
- **Tag** : Tags de transactions

---

*Dernière mise à jour : 29 octobre 2025*
