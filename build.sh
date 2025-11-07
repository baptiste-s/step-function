cat > build.sh << 'EOF'
#!/bin/bash
set -e

# Couleurs pour les messages
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Build Yor Custom Tagging ===${NC}"

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export GOROOT="${SCRIPT_DIR}/go"
export PATH="${GOROOT}/bin:${PATH}"
export GOPATH="${HOME}/go"

# Vérifier Go
echo -e "${YELLOW}Vérification de Go...${NC}"
go version

# Créer le dossier bin
mkdir -p "${SCRIPT_DIR}/bin"

# 1. Compiler Yor
echo -e "${YELLOW}Compilation de Yor...${NC}"
cd "${SCRIPT_DIR}/yor"
go clean
go build -o "${SCRIPT_DIR}/bin/yor" .
echo -e "${GREEN}✓ Yor compilé${NC}"

# 2. Compiler le plugin
echo -e "${YELLOW}Compilation du plugin...${NC}"
cd "${SCRIPT_DIR}/plugin"
go clean
go build -buildmode=plugin -o "${SCRIPT_DIR}/bin/yor_auto_unique_id.so" yor_auto_unique_id.go
echo -e "${GREEN}✓ Plugin compilé${NC}"

echo -e "${GREEN}=== Build terminé ===${NC}"
echo -e "Binaires disponibles dans: ${SCRIPT_DIR}/bin/"
ls -lh "${SCRIPT_DIR}/bin/"
EOF

chmod +x build.sh
