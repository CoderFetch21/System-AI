#!/bin/bash
set -euo pipefail

RED='\033[0;31m' GREEN='\033[0;32m' YELLOW='\033[1;33m' NC='\033[0m'

echo -e "${GREEN}=== SystemAI Installer (github.com/CoderFetch21/System-AI) ===${NC}"

detect_pm() {
    if command -v apt >/dev/null 2>&1; then echo "apt"
    elif command -v pacman >/dev/null 2>&1; then echo "pacman"
    elif command -v dnf >/dev/null 2>&1; then echo "dnf"
    elif command -v zypper >/dev/null 2>&1; then echo "zypper"
    elif command -v emerge >/dev/null 2>&1; then echo "emerge"
    else echo "manual"; fi
}

has_command() { command -v "$1" >/dev/null 2>&1; }

install_dep() {
    local pm="$1" pkg="$2" pkg_cmd="$3"
    echo -e "${YELLOW}Installing $pkg...${NC}"
    echo "Running: sudo -k $pkg_cmd"
    read -p "Continue? (y/N): " -r confirm
    [[ "$confirm" =~ ^[Yy]$ ]] && sudo -k $pkg_cmd && echo -e "${GREEN}✓ Done${NC}"
}

check_deps() {
    local missing=() deps=(go git)
    for dep in "${deps[@]}"; do ! has_command "$dep" && missing+=("$dep"); done
    
    [ ${#missing[@]} -eq 0 ] && { echo -e "${GREEN}✓ Dependencies OK${NC}"; return 0; }
    
    echo -e "${YELLOW}Missing: ${missing[*]}${NC}"
    local pm=$(detect_pm)
    read -p "Auto-install? (y/N): " -r confirm
    
    [[ "$confirm" =~ ^[Yy]$ ]] || { echo -e "${RED}Install manually then retry${NC}"; exit 1; }
    
    case "$pm" in
        apt) install_dep apt go "apt update && apt install -y golang-go git" ;;
        pacman) install_dep pacman go "pacman -Syu --noconfirm go git" ;;
        dnf) install_dep dnf go "dnf install -y golang git" ;;
        zypper) install_dep zypper go "zypper install -y go git" ;;
        emerge) install_dep emerge go "emerge --ask dev-lang/go dev-vcs/git" ;;
        *) echo -e "${RED}Manual install needed${NC}"; exit 1 ;;
    esac
}

build() {
    echo -e "${YELLOW}Building...${NC}"
    mkdir -p build
    go mod tidy
    go build -o build/systemai ./cmd/systemai
    echo -e "${GREEN}✓ Built github.com/CoderFetch21/System-AI${NC}"
}

install_system() {
    echo -e "${YELLOW}System install...${NC}"
    sudo -k cp build/systemai /usr/local/bin/systemai && sudo -k chmod +x /usr/local/bin/systemai
    
    cat > /tmp/systemai.desktop << 'EOF'
[Desktop Entry]
Name=SystemAI
Comment=AI Linux Assistant (github.com/CoderFetch21/System-AI)
Exec=konsole -e systemai
Icon=utilities-terminal
Terminal=false
Type=Application
Categories=System;Utility;
EOF
    sudo -k cp /tmp/systemai.desktop /usr/share/applications/
    echo -e "${GREEN}✓ Installed! Run: systemai or menu${NC}"
}

install_user() {
    local bin="$HOME/.local/bin" apps="$HOME/.local/share/applications"
    mkdir -p "$bin" "$apps"
    cp build/systemai "$bin/systemai" && chmod +x "$bin/systemai"
    
    cat > "$apps/systemai.desktop" << 'EOF'
[Desktop Entry]
Name=SystemAI
Comment=AI Linux Assistant (github.com/CoderFetch21/System-AI)
Exec=konsole -e ~/.local/bin/systemai
Icon=utilities-terminal
Terminal=false
Type=Application
Categories=System;Utility;
EOF
    echo -e "${GREEN}✓ User install! Run: ~/.local/bin/systemai${NC}"
}

main() {
    check_deps && build
    echo -e "\n${YELLOW}Install:${NC} 1=System 2=User"
    read -p "Choice [1]: " -r choice
    case "${choice:-1}" in 1) install_system ;; 2) install_user ;; *) exit 1 ;; esac
}

main "$@"
