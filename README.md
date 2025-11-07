# Yor Custom Tagging

Plugin custom pour Yor permettant d'ajouter un tag unique `carma-name` aux ressources Terraform.

## Prérequis

- **Go 1.19.x** (version exacte requise pour la compatibilité du plugin)

## Installation

### Option 1 : Go 1.19 dans le PATH

Si vous avez déjà Go 1.19.x dans votre PATH :
```bash
git clone <votre-repo-url>
cd yor-custom-tagging

# Vérifier votre version de Go
go version  # Doit afficher go1.19.x

# Compiler
./build.sh
```

### Option 2 : Go 1.19 installé ailleurs

Si Go 1.19 est installé dans un autre emplacement :
```bash
# Compiler en spécifiant le répertoire de Go
./build.sh /chemin/vers/go1.19

# Exemples
./build.sh ~/go1.19
./build.sh /opt/go
./build.sh /data/users/ut3z1p/yor_poc/go
```

**Note:** Le script utilisera le Go spécifié **sans modifier votre PATH**.

### Option 3 : Télécharger Go 1.19

Si vous n'avez pas Go 1.19 :
```bash
# Télécharger
wget https://go.dev/dl/go1.19.13.linux-amd64.tar.gz
tar -xzf go1.19.13.linux-amd64.tar.gz

# Compiler avec ce Go
./build.sh ./go
```

## Utilisation
```bash
# Définir les variables d'environnement
export YOR_ENV="dev"
export YOR_TEAM="myteam"
export YOR_SIMPLE_TAGS='{"Environment": "Dev", "Owner": "Platform"}'

# Exécuter Yor avec le plugin
./bin/yor tag -d ./terraform \
    --parsers Terraform \
    --custom-tagging ./bin/yor_auto_unique_id.so \
    --tag-groups simple

# Dry-run pour prévisualiser
./bin/yor tag -d ./terraform \
    --parsers Terraform \
    --custom-tagging ./bin/yor_auto_unique_id.so \
    --dry-run
```

## Structure du projet
```
yor-custom-tagging/
├── README.md
├── .gitignore
├── build.sh                    # Script de compilation du plugin
├── plugin/
│   ├── yor_auto_unique_id.go  # Code source du plugin
│   ├── go.mod
│   └── go.sum
├── bin/
│   ├── yor                     # Binaire Yor (versionné)
│   └── yor_auto_unique_id.so   # Plugin compilé (généré)
└── examples/
    └── run_yor.sh              # Script d'exemple
```

## Tags générés

### Tag custom `carma-name`

Format : `{env}-{team}-{md5hash}`

Exemple : `dev-myteam-a3f5e9b2c1d4e8f7`

Le hash MD5 est calculé à partir du ResourceID de la ressource Terraform.

### Tags simples (YOR_SIMPLE_TAGS)

Ajoutez `--tag-groups simple` pour activer les tags définis dans `YOR_SIMPLE_TAGS`.

## Développement

### Modifier le plugin

1. Éditez `plugin/yor_auto_unique_id.go`
2. Recompilez :
```bash
./build.sh
# Ou avec un Go spécifique
./build.sh /path/to/go1.19
```

3. Testez :
```bash
export YOR_ENV="test"
export YOR_TEAM="dev"
./bin/yor tag -d ./test-terraform --dry-run --custom-tagging ./bin/yor_auto_unique_id.so
```

### Mettre à jour Yor

Pour mettre à jour le binaire Yor :
```bash
# Compiler Yor depuis les sources avec Go 1.19
cd /tmp
git clone https://github.com/bridgecrewio/yor.git
cd yor
git checkout v0.1.XXX

# Compiler avec Go 1.19
/path/to/go1.19/bin/go build -o yor .

# Copier dans votre repo
cp yor /path/to/yor-custom-tagging/bin/yor

# Recompiler le plugin avec le même Go
cd /path/to/yor-custom-tagging
./build.sh /path/to/go1.19

# Versionner
git add bin/yor
git commit -m "Update Yor to v0.1.XXX"
```

## Troubleshooting

### ❌ Erreur : "plugin was built with a different version"

Le plugin et le binaire Yor doivent être compilés avec **exactement** la même version de Go.

**Solution :**
```bash
# Recompiler le plugin avec Go 1.19
./build.sh /path/to/go1.19
```

### ❌ Le plugin ne charge pas

Vérifiez les logs pour voir si le plugin est bien initialisé :
```bash
./bin/yor tag -d . --custom-tagging ./bin/yor_auto_unique_id.so 2>&1 | grep PLUGIN
```

Vous devriez voir :
```
[PLUGIN] =========== Plugin UniqueID initialisé ===========
[PLUGIN] Initialisation de ExtraTagGroups
[PLUGIN] ExtraTagGroups configuré avec 1 groupe(s)
[PLUGIN] InitTagGroup appelé
[PLUGIN] GetDefaultTags appelé
[PLUGIN] CreateTagsForBlock appelé pour: ...
[PLUGIN] CalculateValue appelé
[PLUGIN] Tag généré: carma-name = dev-myteam-a3f5e9b2c1d4e8f7
```

### ❌ Variables d'environnement non prises en compte

Vérifiez qu'elles sont bien exportées avant d'exécuter Yor :
```bash
echo $YOR_ENV
echo $YOR_TEAM
echo $YOR_SIMPLE_TAGS
```

## Exemples d'utilisation

### Exemple 1 : Tag avec dry-run
```bash
export YOR_ENV="prod"
export YOR_TEAM="platform"

./bin/yor tag -d ./infrastructure \
    --parsers Terraform \
    --custom-tagging ./bin/yor_auto_unique_id.so \
    --dry-run
```

### Exemple 2 : Avec tags simples
```bash
export YOR_ENV="staging"
export YOR_TEAM="backend"
export YOR_SIMPLE_TAGS='{"CostCenter": "Engineering", "Project": "API"}'

./bin/yor tag -d ./terraform \
    --parsers Terraform \
    --custom-tagging ./bin/yor_auto_unique_id.so \
    --tag-groups simple
```

### Exemple 3 : Exclure certains tags par défaut
```bash
export YOR_ENV="dev"
export YOR_TEAM="devops"

./bin/yor tag -d ./terraform \
    --parsers Terraform \
    --custom-tagging ./bin/yor_auto_unique_id.so \
    --skip-tags yor_trace,git_modifiers
```

## Licence

[Votre licence]
EOF
