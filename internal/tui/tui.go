package tui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/CoderFetch21/System-AI/internal/ai"
	"github.com/CoderFetch21/System-AI/internal/config"
	"github.com/CoderFetch21/System-AI/internal/fs"
	"github.com/CoderFetch21/System-AI/internal/pm"
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
		case "exit", "quit":
			return nil
		case "help":
			fmt.Println("Natural language → AI → Execute")
			fmt.Println("Examples: 'update my system', 'install htop', 'show /etc/fstab'")
		case "show config":
			fmt.Printf("%+v\n", cfg)
		default:
			// AI INTERPRETATION FIRST
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
				fmt.Printf("\n❌ Unsafe plan rejected: %v\n", err)
				continue
			}
			
			// SHOW AI PLAN
			fmt.Printf("\n🤖 AI Plan (%d actions):\n", len(plan.Actions))
			if plan.Explanation != "" {
				fmt.Println(plan.Explanation)
			}
			for i, action := range plan.Actions {
				fmt.Printf("  %d. %s", i+1, action.Type)
				if action.Package != "" {
					fmt.Printf(" [%s]", action.Package)
				}
				if action.Path != "" {
					fmt.Printf(" %s", action.Path)
				}
				if action.NeedsRoot {
					fmt.Print(" 🔒")
				}
				fmt.Println()
			}
			
			// CONFIRM EXECUTION
			fmt.Print("\nExecute AI plan? (y/N): ")
			if !confirm(reader) {
				fmt.Println("Plan cancelled.")
				continue
			}
			
			fmt.Println("🚀 Executing AI plan...")
			for i, action := range plan.Actions {
				fmt.Printf("\n--- Action %d/%d ---\n", i+1, len(plan.Actions))
				
				switch action.Type {
				case ai.InstallPackage:
					cmd := pm.InstallCommand(pm.Manager(cfg.PackageManager), action.Package)
					if cmd != nil {
						fmt.Printf("🔄 sudo -k %s %s\n", cfg.PackageManager, action.Package)
						fmt.Println("  (AI would execute this - runner.RunCommand() pending)")
					} else {
						fmt.Printf("❌ No command for %s\n", cfg.PackageManager)
					}
					
				case ai.RunCommand:
					if len(action.Command) > 0 {
						fmt.Printf("🔄 %s\n", strings.Join(action.Command, " "))
						fmt.Println("  (AI would execute this - runner.RunCommand() pending)")
					}
					
				case ai.ReadFile:
					fmt.Printf("📄 Would read: %s\n", action.Path)
					
				default:
					fmt.Printf("⚠️ %s action pending implementation\n", action.Type)
				}
			}
			fmt.Println("\n✅ AI plan processed!")
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
