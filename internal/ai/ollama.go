package ai

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "github.com/yourusername/systemai/internal/config"
)

type OllamaPlanner struct {
    model string
    endpoint string
}

func NewOllamaPlanner(cfg *config.Config) *OllamaPlanner {
    return &OllamaPlanner{
        model: cfg.AiModel,
        endpoint: "http://127.0.0.1:11434/api/generate",
    }
}

func (p *OllamaPlanner) Plan(input string, ctx Context) (*Plan, error) {
    prompt := buildSystemPrompt(ctx)
    prompt += fmt.Sprintf("\n\nUser: %s\n\n
