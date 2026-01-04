package pm

import (
	"os/exec"
	"strings"
)

type Manager string

const (
	Manual  Manager = "manual"
	Apt     Manager = "apt"
	Pacman  Manager = "pacman"
	Dnf     Manager = "dnf"
	Zypper  Manager = "zypper"
	Emerge  Manager = "emerge"
)

func Detect() Manager {
	if hasCommand("apt") { return Apt }
	if hasCommand("pacman") { return Pacman }
	if hasCommand("dnf") { return Dnf }
	if hasCommand("zypper") { return Zypper }
	if hasCommand("emerge") { return Emerge }
	return Manual
}

func InstallCommand(m Manager, pkg string) []string {
	sudo := []string{"sudo", "-k"}
	switch m {
	case Apt:    return append(sudo, "apt", "install", "-y", pkg)
	case Pacman: return append(sudo, "pacman", "-S", "--noconfirm", pkg)
	case Dnf:    return append(sudo, "dnf", "install", "-y", pkg)
	case Zypper: return append(sudo, "zypper", "install", "-y", pkg)
	case Emerge: return append(sudo, "emerge", "--ask", pkg)
	default:     return nil
	}
}

func UpdateSystemCommands(m Manager) [][]string {
	switch m {
	case Emerge:
		return [][]string{
			{"sudo", "-k", "emaint", "sync", "-a"},
			{"sudo", "-k", "emerge", "--ask", "--update", "--deep", "--newuse", "@world"},
		}
	case Apt:
		return [][]string{{"sudo", "-k", "apt", "update", "-y", "&&", "apt", "upgrade", "-y"}}
	case Pacman:
		return [][]string{{"sudo", "-k", "pacman", "-Syu"}}
	default:
		return nil
	}
}

func hasCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}
