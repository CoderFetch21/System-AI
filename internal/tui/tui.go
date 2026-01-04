package tui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/CoderFetch21/System-AI/internal/ai"
	"github.com/CoderFetch21/System-AI/internal/config"
	"github.com/CoderFetch21/System-AI/internal/fs"
	"github.com/CoderFetch21/System-AI/internal/pm"
	"github.com/CoderFetch21/System-AI/internal/runner"
)

func RunFirstRunWizard() (*config.Config, error) {
	cfg := &config.Config{AiBackend: "ollama", AiModel: "llama3.2:3b"}
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
	fmt.Println("🧠 SystemAI + Llama 3.2 3B Ready!")
	fmt.Println("Try: 'update my system', 'install htop', 'show /etc/fstab'")
	
	for {
		fmt.Print("systemai> ")
		input := readLine(reader)
		
		switch input {
		case "exit", "quit": return nil
		case "help":
			fmt.Println("Natural language → AI → Execute")
		case "show config":
			fmt.Printf("%+v\n", cfg)
		default:
			// AI FIRST, THEN EXECUTE
			planner := ai.NewOllamaPlanner(cfg)
			aiCtx := ai.Context{
				DistroFamily:   cfg.DistroFamily,
				PackageManager: cfg.PackageManager,
				Cwd:            "/",
				UserQuery:      input,
			}
			
			fmt.Print("🧠 AI interpreting...")
			plan, err := planner.Plan(aiCtx)
			if err != nil {
				fmt.Printf("\n❌ AI error: %v\n", err)
				continue
			}
			
			if err := planner.Validate(plan); err != nil {
				fmt.Printf("\n❌ Unsafe plan: %v\n", err)
				continue
			}
			
			// SHOW AI PLAN
			fmt.Printf("\n🤖 AI Plan (%d actions):\n", len(plan.Actions))
			fmt.Println(plan.Explanation)
			for i, action := range plan.Actions {
				fmt.Printf("  %d. %s", i+1, action.Type)
				if action.Package != "" { fmt.Printf(" [%s]", action.Package) }
				if action.Path != "" { fmt.Printf(" %s", action.Path) }
				if action.NeedsRoot { fmt.Print(" 🔒") }
				fmt.Println()
			}
			
			// EXECUTE
			fmt.Print("\nExecute AI plan? (y/N): ")
			if !confirm(reader) {
				fmt.Println("Cancelled.")
				continue
			}
			
			fmt.Println("🚀 Executing...")
			for i, action := range plan.Actions {
				fmt.Printf("\n--- Action %d/%d ---\n", i+1, len(plan.Actions))
				
				switch action.Type {
				case ai.InstallPackage:
					cmd := pm.InstallCommand(pm.Manager(cfg.PackageManager), action.Package)
					if cmd != nil && runner.RunCommand(cmd) == nil {
						fmt.Println("✅ Installed")
					}
					
				case ai.RunCommand:
					if runner.RunCommand(action.Command) == nil {
						fmt.Println("✅ Command OK")
					}
					
				case ai.ReadFile:
					data, err := fs.ReadFileUser(action.Path)
					if fs.IsPermissionError(err) && cfg.AllowRootExecute {
						fmt.Print("Root access needed. Retry? (y/N): ")
						if confirm(reader) {
							data, err = fs.ReadFileRoot(action.Path)
						}
					}
					if err != nil {
						fmt.Printf("❌ Read error: %v\n", err)
					} else {
						fmt.Printf("📄 %s:\n%s\n", action.Path, truncate(string(data), 500))
					}
					
				default:
					fmt.Printf("⚠️ %s pending\n", action.Type)
				}
			}
			fmt.Println("\n✅ AI plan complete!")
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

func truncate(s string, max int) string {
	if len(s) <= max { return s }
	return s[:max] + "..."
}
