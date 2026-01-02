package config

import (
    "encoding/json"
    "os"
)

type Config struct {
    FirstRunCompleted bool   `json:"first_run_completed"`
    DistroFamily      string `json:"distro_family"`
    PackageManager    string `json:"package_manager"`
    Shell             string `json:"shell"`
    Editor            string `json:"editor"`
    AllowRootSuggest  bool   `json:"allow_root_suggest"`
    AllowRootExecute  bool   `json:"allow_root_execute"`
    AiBackend         string `json:"ai_backend"`
    AiModel           string `json:"ai_model"`
}

func Load(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}

func Save(path string, cfg *Config) error {
    cfg.FirstRunCompleted = true
    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(path, data, 0o644)
}
