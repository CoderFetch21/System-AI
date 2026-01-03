#!/bin/bash
set -euo pipefail

RED='\033[0;31m' GREEN='\033[0;32m' YELLOW='\033[1;33m' NC='\033[0m'

echo -e "${GREEN}=== SystemAI Uninstaller ===${NC}"

# Ask how it was installed
echo "How was SystemAI installed?"
echo "1) System-wide (via install.sh to /usr/local/bin and /usr/share/applications)"
echo "2) User-local (to ~/.local/bin and ~/.local/share/applications)"
read -p "Select [1]: " -r choice
choice="${choice:-1}"

uninstall_system() {
    echo -e "${YELLOW}Removing system-wide installation...${NC}"

    # Remove binary
    if command -v systemai >/dev/null 2>&1; then
        bin_path="$(command -v systemai)"
        if [[ "$bin_path" == "/usr/local/bin/systemai" ]]; then
            echo "About to run: sudo -k rm -f /usr/local/bin/systemai"
            read -p "Continue? (y/N): " -r c
            if [[ "$c" =~ ^[Yy]$ ]]; then
                sudo -k rm -f /usr/local/bin/systemai
                echo -e "${GREEN}✓ Removed /usr/local/bin/systemai${NC}"
            fi
        else
            echo -e "${YELLOW}systemai found at $bin_path (not /usr/local/bin), leaving it.${NC}"
        fi
    else
        echo "systemai binary not found in PATH."
    fi

    # Remove desktop file
    if [[ -f /usr/share/applications/systemai.desktop ]]; then
        echo "About to run: sudo -k rm -f /usr/share/applications/systemai.desktop"
        read -p "Continue? (y/N): " -r c
        if [[ "$c" =~ ^[Yy]$ ]]; then
            sudo -k rm -f /usr/share/applications/systemai.desktop
            echo -e "${GREEN}✓ Removed /usr/share/applications/systemai.desktop${NC}"
        fi
    else
        echo "No /usr/share/applications/systemai.desktop found."
    fi

    echo -e "${YELLOW}User configuration in ~/.config/systemai is left intact.${NC}"
}

uninstall_user() {
    echo -e "${YELLOW}Removing user-local installation...${NC}"

    # Remove binary
    if [[ -f "$HOME/.local/bin/systemai" ]]; then
        echo "Removing $HOME/.local/bin/systemai"
        rm -f "$HOME/.local/bin/systemai"
        echo -e "${GREEN}✓ Removed ~/.local/bin/systemai${NC}"
    else
        echo "No ~/.local/bin/systemai found."
    fi

    # Remove desktop file
    if [[ -f "$HOME/.local/share/applications/systemai.desktop" ]]; then
        echo "Removing $HOME/.local/share/applications/systemai.desktop"
        rm -f "$HOME/.local/share/applications/systemai.desktop"
        echo -e "${GREEN}✓ Removed user .desktop file${NC}"
    else
        echo "No ~/.local/share/applications/systemai.desktop found."
    fi

    echo -e "${YELLOW}User configuration in ~/.config/systemai is left intact.${NC}"
}

case "$choice" in
    1) uninstall_system ;;
    2) uninstall_user ;;
    *)
        echo -e "${RED}Invalid choice, aborting.${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}SystemAI uninstallation finished.${NC}"
echo "If you also want to remove config, delete: ~/.config/systemai"
