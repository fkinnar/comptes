# Mise à jour de Go pour macOS 26+

Sur macOS 26+ (Sequoia), Go 1.22.0 et versions antérieures ne génèrent pas le `LC_UUID` requis par `dyld`, ce qui cause l'erreur :
```
dyld: missing LC_UUID load command
```

## Solution : Mettre à jour Go

### Méthode 1 : Via Homebrew (recommandé)
```bash
# Installer/mettre à jour Go vers la dernière version
brew install go
# ou
brew upgrade go

# Vérifier la version (doit être 1.24+)
go version

# Si Homebrew Go est installé mais pas dans le PATH

## Pour Zsh/Bash
Ajoutez ceci à votre ~/.zshrc ou ~/.bashrc:
```bash
export PATH="/opt/homebrew/bin:$PATH"
```

## Pour NuShell
Ajoutez ceci à votre ~/.config/nushell/config.nu:
```nu
# Ajouter /opt/homebrew/bin au PATH
if "/opt/homebrew/bin" not-in $env.PATH {
    $env.PATH = ($env.PATH | prepend "/opt/homebrew/bin")
}
```

## Ou utilisez directement:
```bash
/opt/homebrew/bin/go version
```

### Méthode 2 : Script automatique
```bash
# Utiliser le script fourni
./update-go.sh
```

### Méthode 3 : Installation manuelle
1. Téléchargez Go 1.24+ depuis https://golang.org/dl/
2. Installez le package
3. Vérifiez avec `go version`

## Après la mise à jour

Une fois Go mis à jour, vous pouvez compiler normalement :
```bash
make build
./comptes list
```

## Vérification

Pour vérifier que le problème est résolu :
```bash
go version  # Doit afficher go1.24.x ou supérieur
make clean
make build
./comptes list  # Ne devrait plus afficher l'erreur LC_UUID
```

