#!/bin/bash
set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== Setup Yor Custom Tagging ===${NC}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Vérifier que Go est présent
if [ ! -d "${SCRIPT_DIR}/go" ]; then
    echo -e "${RED}Erreur: Go n'est pas présent dans ${SCRIPT_DIR}/go${NC}"
    echo "Copiez votre installation Go 1.19 dans ce répertoire"
    exit 1
fi

# Init submodule si nécessaire
if [ -d "${SCRIPT_DIR}/yor/.git" ]; then
    echo -e "${YELLOW}Mise à jour du submodule Yor...${NC}"
    git submodule update --init --recursive
fi

# Vérifier les dépendances du plugin
echo -e "${YELLOW}Vérification des dépendances du plugin...${NC}"
cd "${SCRIPT_DIR}/plugin"
export GOROOT="${SCRIPT_DIR}/go"
export PATH="${GOROOT}/bin:${PATH}"

if [ ! -f "go.mod" ]; then
    echo -e "${YELLOW}Initialisation du module Go...${NC}"
    go mod init yor_auto_unique_id
    go get github.com/bridgecrewio/yor@v0.1.200
fi

go mod tidy

echo -e "${GREEN}✓ Setup terminé${NC}"
echo -e "Lancez ${YELLOW}./build.sh${NC} pour compiler"
EOF
