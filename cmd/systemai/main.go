package main

import (
    "fmt"
    "log"
    "os"
    "os/user"
    "path/filepath"

    "github.com/yourusername/systemai/internal/config"
    "github.com/yourusername/systemai/internal/tui"
)

func main() {
    usr, err := user.Current()
    if err != nil {
        log.Fatalf("failed to get current user: %v", err)
    }

    configDir := filepath.Join(usr.HomeDir, ".config", "systemai")
    if err := os.MkdirAll(configDir, 0o755); err != nil {
        log.Fatalf("failed to create config dir: %v", err)
    }

    configPath := filepath.Join(configDir, "config.json")
    cfg, err := config.Load(configPath)
    if err != nil {
        fmt.Println("🧠 SystemAI First-Run Setup")
        cfg, err = tui.RunFirstRunWizard()
        if err != nil {
            log.Fatalf("first-run wizard failed: %v", err)
        }
        if err := config.Save(configPath, cfg); err != nil {
            log.Fatalf("failed to save config: %v", err)
        }
    }

    if err := tui.RunMainTUI(cfg, configPath); err != nil {
        log.Fatalf("SystemAI failed: %v", err)
    }
}
