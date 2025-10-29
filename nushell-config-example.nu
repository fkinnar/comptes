# Configuration NuShell pour Comptes
# Ajoute /opt/homebrew/bin au PATH pour utiliser Go 1.25.3+

# Vérifier si /opt/homebrew/bin n'est pas déjà dans le PATH
if "/opt/homebrew/bin" not-in $env.PATH {
    $env.PATH = ($env.PATH | prepend "/opt/homebrew/bin")
}

# Alternative: utiliser la syntaxe liste
# $env.PATH = [
#     "/opt/homebrew/bin"
#     ...$env.PATH
# ]

