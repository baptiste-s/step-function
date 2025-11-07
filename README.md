# Yor Custom Tagging

Plugin custom pour Yor permettant d'ajouter un tag unique `carma-name` aux ressources Terraform, ainsi que des tags configurables via YAML.

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
├── yaml/
│   └── tags.yaml               # Configuration des tags YAML
├── bin/
│   ├── yor                     # Binaire Yor (versionné)
│   └── yor_auto_unique_id.so   # Plugin compilé (généré)
└── examples/
    └── run_yor.sh              # Script d'exemple
```

## Tags disponibles

### 1. Tag custom via plugin : `carma-name`

**Format :** `{env}-{team}-{md5hash}`

**Exemple :** `dev-myteam-a3f5e9b2c1d4e8f7`

Le hash MD5 est calculé à partir du ResourceID de la ressource Terraform.

**Variables d'environnement requises :**
- `YOR_ENV` : Environnement (dev, staging, prod, etc.)
- `YOR_TEAM` : Équipe propriétaire

### 2. Tags CA via YAML : `ca-applicationid`, `ca-owner`, `ca-company`

Définis dans `yaml/tags.yaml` :

- **ca-applicationid** : ID de l'application
- **ca-owner** : Propriétaire/équipe responsable
- **ca-company** : Nom de l'entreprise

**Variables d'environnement requises :**
- `CA_APPLICATION_ID`
- `CA_OWNER`
- `CA_COMPANY`

### 3. Tags simples (YOR_SIMPLE_TAGS)

Tags clé-valeur statiques définis via la variable d'environnement `YOR_SIMPLE_TAGS`.

## Utilisation

### Option 1 : Plugin custom uniquement
```bash
export YOR_ENV="dev"
export YOR_TEAM="myteam"

./bin/yor tag -d ./terraform \
    --parsers Terraform \
    --custom-tagging ./bin/yor_auto_unique_id.so
```

### Option 2 : Tags YAML uniquement
```bash
export CA_APPLICATION_ID="APP-12345"
export CA_OWNER="team-platform"
export CA_COMPANY="CACIB"

./bin/yor tag -d ./terraform \
    --parsers Terraform \
    --config-file ./yaml/tags.yaml
```

### Option 3 : Plugin + YAML (recommandé)
```bash
# Variables pour le plugin custom
export YOR_ENV="dev"
export YOR_TEAM="myteam"

# Variables pour les tags YAML
export CA_APPLICATION_ID="APP-12345"
export CA_OWNER="team-platform"
export CA_COMPANY="CACIB"

./bin/yor tag -d ./terraform \
    --parsers Terraform \
    --config-file ./yaml/tags.yaml \
    --custom-tagging ./bin/yor_auto_unique_id.so
```

### Option 4 : Tout combiné (Plugin + YAML + Simple)
```bash
# Variables pour le plugin custom
export YOR_ENV="prod"
export YOR_TEAM="platform"

# Variables pour les tags YAML
export CA_APPLICATION_ID="APP-12345"
export CA_OWNER="team-platform"
export CA_COMPANY="CACIB"

# Variables pour les tags simples
export YOR_SIMPLE_TAGS='{"Environment": "Production", "CostCenter": "Engineering"}'

./bin/yor tag -d ./terraform \
    --parsers Terraform \
    --config-file ./yaml/tags.yaml \
    --custom-tagging ./bin/yor_auto_unique_id.so \
    --tag-groups simple
```

**Résultat :** Chaque ressource aura tous ces tags :
- `carma-name` (plugin)
- `ca-applicationid`, `ca-owner`, `ca-company` (YAML)
- `Environment`, `CostCenter` (simple tags)
- Tags par défaut de Yor : `yor_trace`, `git_*`, etc.

### Dry-run (prévisualisation sans modification)
```bash
./bin/yor tag -d ./terraform \
    --parsers Terraform \
    --config-file ./yaml/tags.yaml \
    --custom-tagging ./bin/yor_auto_unique_id.so \
    --dry-run
```

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

### Modifier les tags YAML

Éditez `yaml/tags.yaml` pour ajouter ou modifier des tags. Aucune recompilation nécessaire.

**Exemple - Ajouter un nouveau tag :**
```yaml
tag_groups:
  - name: ca_tags
    tags:
      - name: ca-applicationid
        value:
          default: ${env:CA_APPLICATION_ID}
      
      - name: ca-owner
        value:
          default: ${env:CA_OWNER}
      
      - name: ca-company
        value:
          default: ${env:CA_COMPANY}
      
      # Nouveau tag
      - name: ca-project
        value:
          default: ${env:CA_PROJECT}
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

### ❌ Tags YAML non appliqués

Vérifiez que :
1. Les variables d'environnement sont bien définies :
```bash
echo $CA_APPLICATION_ID
echo $CA_OWNER
echo $CA_COMPANY
```

2. Le path du fichier YAML est correct :
```bash
ls -la ./yaml/tags.yaml
```

3. Le fichier YAML est valide :
```bash
# Vérifier la syntaxe YAML (optionnel)
python3 -c "import yaml; yaml.safe_load(open('./yaml/tags.yaml'))"
```

### ❌ Variables d'environnement non prises en compte

Si les variables d'environnement ne sont pas prises en compte, vérifiez qu'elles sont **exportées** :
```bash
# ✅ Correct
export YOR_ENV="dev"

# ❌ Incorrect (pas exporté)
YOR_ENV="dev"
```

## Exemples d'utilisation

### Exemple 1 : Tag complet pour production
```bash
export YOR_ENV="prod"
export YOR_TEAM="platform"
export CA_APPLICATION_ID="APP-PROD-001"
export CA_OWNER="team-platform"
export CA_COMPANY="CACIB"
export YOR_SIMPLE_TAGS='{"Environment": "Production", "Compliance": "SOC2"}'

./bin/yor tag -d ./infrastructure \
    --parsers Terraform \
    --config-file ./yaml/tags.yaml \
    --custom-tagging ./bin/yor_auto_unique_id.so \
    --tag-groups simple
```

### Exemple 2 : Tag avec dry-run pour validation
```bash
export YOR_ENV="staging"
export YOR_TEAM="backend"
export CA_APPLICATION_ID="APP-STG-002"
export CA_OWNER="team-backend"
export CA_COMPANY="CACIB"

./bin/yor tag -d ./terraform \
    --parsers Terraform \
    --config-file ./yaml/tags.yaml \
    --custom-tagging ./bin/yor_auto_unique_id.so \
    --dry-run
```

### Exemple 3 : Exclure certains tags par défaut
```bash
export YOR_ENV="dev"
export YOR_TEAM="devops"
export CA_APPLICATION_ID="APP-DEV-003"
export CA_OWNER="team-devops"
export CA_COMPANY="CACIB"

./bin/yor tag -d ./terraform \
    --parsers Terraform \
    --config-file ./yaml/tags.yaml \
    --custom-tagging ./bin/yor_auto_unique_id.so \
    --skip-tags yor_trace,git_modifiers
```

### Exemple 4 : Tag uniquement des fichiers spécifiques
```bash
export YOR_ENV="dev"
export YOR_TEAM="data"
export CA_APPLICATION_ID="APP-DEV-004"
export CA_OWNER="team-data"
export CA_COMPANY="CACIB"

# Tagger uniquement le répertoire s3
./bin/yor tag -d ./terraform/s3 \
    --parsers Terraform \
    --config-file ./yaml/tags.yaml \
    --custom-tagging ./bin/yor_auto_unique_id.so
```

## Récapitulatif des tags appliqués

Avec la configuration complète (Plugin + YAML + Simple + Défaut), voici tous les tags appliqués :

| Tag | Source | Exemple de valeur |
|-----|--------|-------------------|
| `carma-name` | Plugin custom | `dev-myteam-a3f5e9b2c1d4e8f7` |
| `ca-applicationid` | YAML | `APP-12345` |
| `ca-owner` | YAML | `team-platform` |
| `ca-company` | YAML | `CACIB` |
| `Environment` | YOR_SIMPLE_TAGS | `Production` |
| `CostCenter` | YOR_SIMPLE_TAGS | `Engineering` |
| `yor_trace` | Yor défaut | `912066a1-31a3-4a08-911b-0b06d9eac64e` |
| `git_repo` | Yor défaut | `yor-custom-tagging` |
| `git_org` | Yor défaut | `myorg` |
| `git_file` | Yor défaut | `terraform/main.tf` |
| `git_commit` | Yor défaut | `a1b2c3d4...` |
| `git_modifiers` | Yor défaut | `johndoe/janedoe` |
| `git_last_modified_at` | Yor défaut | `2025-01-15 10:30:00` |
| `git_last_modified_by` | Yor défaut | `johndoe@company.com` |

## Licence


