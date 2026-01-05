package ai

import (
	"bytes"
	"encoding/json"
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
	model      string
	endpoint   string
	httpClient *http.Client
}

func NewOllamaPlanner(cfg *config.Config) *OllamaPlanner {
	model := cfg.AiModel
	if model == "" {
		model = "llama3.2:3b"
	}
	return &OllamaPlanner{
		model:    model,
		endpoint: "http://127.0.0.1:11434/api/generate",
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (p *OllamaPlanner) Plan(ctx Context) (*Plan, error) {
	prompt := p.buildPrompt(ctx)

	reqBody := map[string]any{
		"model":  p.model,
		"prompt": prompt,
		"stream": false,
		"format": "json",
		"options": map[string]any{
			"temperature": 0.1,
			"top_p":       0.9,
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := p.httpClient.Post(p.endpoint, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama HTTP %d: %s", resp.StatusCode, string(b))
	}

	var ollamaResp struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("decode ollama response: %w", err)
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
	return fmt.Sprintf(`You are SystemAI, a Linux system assistant that plans actions but never executes them directly.

Context:
- Distro: %s
- Package manager: %s
- Current directory: %s

User request:
%s

You MUST output ONLY valid JSON with this schema:

{
  "actions": [
    {
      "type": "install_package" | "remove_package" | "read_file" | "edit_file" | "create_file" | "run_command",
      "package": "name or empty",
      "path": "/absolute/path or empty",
      "language": "bash|nginx|json|yaml|other or empty",
      "content": "full file content when creating a file",
      "diff": "unified diff when editing a file",
      "command": ["cmd", "arg1", "arg2"],
      "needs_root": true or false
    }
  ],
  "explanation": "brief natural language explanation of the plan"
}

Rules:
- Use ONLY the action types listed.
- For Gentoo with emerge, prefer commands like:
  - "sudo -k emaint sync -a"
  - "sudo -k emerge --ask --update --deep --newuse @world"
- For install_package, set "package" to the package name only; SystemAI will map it to the right command.
- For run_command, fill "command" as an array of tokens.
- Set needs_root=true for system-level operations (e.g., /etc, package management, systemctl, etc.).
- Prefer minimal, safe, and reversible changes (use diffs for config edits).
- DO NOT include any text outside the JSON object.`,
		ctx.DistroFamily,
		ctx.PackageManager,
		ctx.Cwd,
		ctx.UserQuery,
	)
}

func (p *OllamaPlanner) Validate(plan *Plan) error {
	dangerousPaths := []string{"/", "/boot", "/proc", "/sys", "/dev"}

	for i, action := range plan.Actions {
		switch action.Type {
		case InstallPackage, RemovePackage:
			if action.Package == "" {
				return fmt.Errorf("action %d: empty package name", i)
			}
		case ReadFile, EditFile, CreateFile:
			if action.Path == "" {
				return fmt.Errorf("action %d: empty path", i)
			}
			for _, d := range dangerousPaths {
				if action.Path == d {
					return fmt.Errorf("action %d: path %s is too dangerous", i, action.Path)
				}
			}
		case RunCommand:
			if len(action.Command) == 0 {
				return fmt.Errorf("action %d: empty command", i)
			}
			cmdStr := strings.Join(action.Command, " ")
			if strings.Contains(cmdStr, "rm -rf /") || strings.Contains(cmdStr, "mkfs") {
				return fmt.Errorf("action %d: dangerous command: %s", i, cmdStr)
			}
		default:
			return fmt.Errorf("action %d: unknown type %q", i, action.Type)
		}
	}
	return nil
}
