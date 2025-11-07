# Yor Custom Tagging

Plugin custom pour Yor permettant d'ajouter des tags uniques aux ressources Terraform.

## Prérequis

- Go 1.19 
- Git

## Installation
```bash
# Cloner le repo
git clone <votre-repo-url>
cd yor-custom-tagging

# Setup initial
./setup.sh

# Compiler
./build.sh
```

## Utilisation
```bash
# Variables d'environnement
export YOR_ENV="dev"
export YOR_TEAM="myteam"

# Exécuter Yor avec le plugin
./bin/yor tag -d ./terraform \
    --parsers Terraform \
    --custom-tagging ./bin/yor_auto_unique_id.so \
    --tag-groups simple

# Ou utiliser le script d'exemple
export YOR_ENV="prod"
export YOR_TEAM="platform"
./examples/run_yor.sh ./terraform --dry-run
```

## Structure

- `yor/` - Sources de Yor 
- `plugin/` - Code source du plugin custom
- `bin/` - Binaires compilés (non versionné)
- `examples/` - Scripts d'exemple

## Développement

### Modifier le plugin

1. Éditez `plugin/yor_auto_unique_id.go`
2. Recompilez : `./build.sh`
3. Testez : `./examples/run_yor.sh ./test-terraform --dry-run`

### Mettre à jour Yor
```bash
cd yor
git checkout <nouvelle-version>
cd ..
git add yor
git commit -m "Update Yor to <version>"
./build.sh  # Recompiler
```

## Tags ajoutés

- `carma-name`: `{env}-{team}-{md5hash}`
- Tags de `YOR_SIMPLE_TAGS` si `--tag-groups simple` est spécifié

## Troubleshooting

### Erreur "plugin was built with a different version"

Assurez-vous que Yor et le plugin sont compilés avec le même Go :
```bash
./build.sh
```

### Le plugin ne charge pas

Vérifiez les logs :
```bash
./bin/yor tag -d . --custom-tagging ./bin/yor_auto_unique_id.so 2>&1 | grep PLUGIN
```
EOF
