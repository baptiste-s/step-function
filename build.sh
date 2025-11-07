#!/bin/bash
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== Build Plugin Custom Yor ===${NC}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Déterminer quel Go utiliser
if [ -n "$1" ]; then
    GO_BIN="$1/bin/go"
    echo -e "${YELLOW}Utilisation de Go depuis: $1${NC}"
    
    if [ ! -f "$GO_BIN" ]; then
        echo -e "${RED}❌ Erreur: Go introuvable dans $1/bin/go${NC}"
        exit 1
    fi
else
    GO_BIN="go"
    
    if ! command -v go &> /dev/null; then
        echo -e "${RED}❌ Erreur: Go n'est pas installé ou introuvable dans le PATH${NC}"
        echo ""
        echo "📥 Pour télécharger Go 1.19.13:"
        echo "   wget https://go.dev/dl/go1.19.13.linux-amd64.tar.gz"
        echo "   tar -xzf go1.19.13.linux-amd64.tar.gz"
        echo ""
        echo "🔧 Puis relancez avec le path de Go:"
        echo "   ./build.sh /path/to/go"
        echo ""
        echo "Exemple:"
        echo "   ./build.sh ~/go1.19"
        exit 1
    fi
fi

# Vérifier la version de Go
GO_VERSION=$($GO_BIN version | awk '{print $3}' | sed 's/go//')
GO_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
GO_MINOR=$(echo $GO_VERSION | cut -d. -f2)

echo -e "${YELLOW}Version de Go détectée: ${GO_VERSION}${NC}"

if [ "$GO_MAJOR" != "1" ] || [ "$GO_MINOR" != "19" ]; then
    echo -e "${RED}❌ Erreur: Go 1.19.x est requis, vous avez Go ${GO_VERSION}${NC}"
    echo ""
    echo "📥 Pour télécharger Go 1.19.13:"
    echo "   wget https://go.dev/dl/go1.19.13.linux-amd64.tar.gz"
    echo "   tar -xzf go1.19.13.linux-amd64.tar.gz"
    echo ""
    echo "🔧 Puis relancez avec le path de Go:"
    echo "   ./build.sh /path/to/go"
    echo ""
    echo "Exemple:"
    echo "   ./build.sh ~/go1.19"
    exit 1
fi

echo -e "${GREEN}✓ Go 1.19.x OK${NC}"

# Créer le dossier bin
mkdir -p "${SCRIPT_DIR}/bin"

# Vérifier les dépendances
echo -e "${YELLOW}Vérification des dépendances...${NC}"
cd "${SCRIPT_DIR}/plugin"

if [ ! -f "go.mod" ]; then
    echo -e "${YELLOW}Initialisation du module Go...${NC}"
    $GO_BIN mod init yor_auto_unique_id
    $GO_BIN get github.com/bridgecrewio/yor@v0.1.200
fi

$GO_BIN mod tidy

# Compiler le plugin
echo -e "${YELLOW}Compilation du plugin...${NC}"
$GO_BIN clean -cache
$GO_BIN build -buildmode=plugin -o "${SCRIPT_DIR}/bin/yor_auto_unique_id.so" yor_auto_unique_id.go

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Plugin compilé avec succès${NC}"
    echo ""
    ls -lh "${SCRIPT_DIR}/bin/yor_auto_unique_id.so"
    echo ""
    echo -e "${GREEN}=== Build terminé ===${NC}"
else
    echo -e "${RED}✗ Erreur lors de la compilation${NC}"
    exit 1
fi

