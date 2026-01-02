package tui

import (
    "bufio"
    "fmt"
    "os"
    "strings"

    "github.com/CoderFetch21/System-AI/internal/config"
)

func RunFirstRunWizard() (*config.Config, error) {
    cfg := &config.Config{
        AiBackend: "ollama",
        AiModel:   "llama3.2:3b",
    }
    reader := bufio.NewReader(os.Stdin)

    fmt.Print("Distro family (debian/arch/fedora/gentoo/other): ")
    cfg.DistroFamily = readLine(reader)
    fmt.Print("Package manager (apt/pacman/dnf/zypper/emerge/manual): ")
    cfg.PackageManager = readLine(reader)
    fmt.Print("Shell (bash/zsh/fish/other): ")
    cfg.Shell = readLine(reader)
    fmt.Print("Editor (nano/vim/micro/other): ")
    cfg.Editor = readLine(reader)
    fmt.Print("Allow root suggestions? (y/N): ")
    cfg.AllowRootSuggest = confirm(reader)
    fmt.Print("Allow root execution? (y/N): ")
    cfg.AllowRootExecute = confirm(reader)

    return cfg, nil
}

func RunMainTUI(cfg *config.Config, configPath string) error {
    reader := bufio.NewReader(os.Stdin)
    fmt.Println("🧠 SystemAI Ready! Type 'help' or natural language commands")
    
    for {
        fmt.Print("systemai> ")
        input := readLine(reader)
        switch input {
        case "exit", "quit": return nil
        case "help":
            fmt.Println("help, exit, show config, or ask anything (install, edit files, etc)")
        case "show config":
            fmt.Printf("Config: %+v\n", cfg)
        default:
            fmt.Printf("🧠 Processing: %q (AI integration WIP)\n", input)
        }
    }
}

func readLine(r *bufio.Reader) string {
    text, _ := r.ReadString('\n')
    return strings.TrimSpace(text)
}

func confirm(r *bufio.Reader) bool {
    return strings.ToLower(readLine(r)) == "y"
}
