package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/CoderFetch21/System-AI/internal/config"
)

type Context struct {
	DistroFamily   string   `json:"distro_family"`
	PackageManager string   `json:"package_manager"`
	Cwd            string   `json:"cwd"`
	UserQuery      string   `json:"user_query"`
	RecentActions  []string `json:"recent_actions,omitempty"`
}

type OllamaPlanner struct {
	model     string
	endpoint  string
	httpClient *http.Client
}

func NewOllamaPlanner(cfg *config.Config) *OllamaPlanner {
	return &OllamaPlanner{
		model:     cfg.AiModel,
		endpoint:  "http://127.0.0.1:11434/api/generate",
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func (p *OllamaPlanner) Plan(ctx Context) (*Plan, error) {
	prompt := p.buildPrompt(ctx)
	
	reqBody := map[string]interface{}{
		"model":  p.model,
		"prompt": prompt,
		"stream": false,
		"format": "json",
		"options": map[string]interface{}{
			"temperature": 0.1,
			"top_p":       0.9,
		},
	}
	
	bodyBytes, _ := json.Marshal(reqBody)
	resp, err := p.httpClient.Post(p.endpoint, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama HTTP %d: %s", resp.StatusCode, string(body))
	}
	
	var ollamaResp struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to parse ollama response: %w", err)
	}
	
	var plan Plan
	if err := json.Unmarshal([]byte(ollamaResp.Response), &plan); err != nil {
		return nil, fmt.Errorf("invalid JSON plan from ollama: %w", err)
	}
	
	if len(plan.Actions) == 0 {
		plan.Explanation = "No actions identified from user request."
	}
	
	return &plan, nil
}

func (p *OllamaPlanner) buildPrompt(ctx Context) string {
	systemPrompt := fmt.Sprintf(`You are SystemAI, a Linux system assistant. 

Context:
- Distro: %s
- Package manager: %s

User: %s

Output ONLY valid JSON:
{
  "actions": [
    {
      "type": "install_package|read_file|run_command",
      "package": "name",
      "path": "/full/path", 
      "command": ["sudo", "-k", "emerge", "--ask", "sys-process/htop"],
      "needs_root": true
    }
  ],
  "explanation": "brief explanation"
}`, ctx.DistroFamily, ctx.PackageManager, ctx.UserQuery)
	
	return systemPrompt
}

func (p *OllamaPlanner) Validate(plan *Plan) error {
	dangerousPaths := []string{"/", "/boot", "/proc", "/sys", "/dev"}
	
	for i, action := range plan.Actions {
		switch action.Type {
		case InstallPackage, RunCommand:
			if len(action.Command) == 0 && action.Package == "" {
				return fmt.Errorf("action %d: empty command/package", i)
			}
		case ReadFile:
			if action.Path == "" {
				return fmt.Errorf("action %d: empty path", i)
			}
			for _, dangerous := range dangerousPaths {
				if strings.HasPrefix(action.Path, dangerous) {
					return fmt.Errorf("action %d: dangerous path %s", i, action.Path)
				}
			}
		}
	}
	return nil
}
