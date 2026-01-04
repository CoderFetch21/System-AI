package fs

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func ReadFileUser(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func ReadFileRoot(path string) ([]byte, error) {
	cmd := exec.Command("sudo", "-k", "cat", path)
	return cmd.CombinedOutput()
}

func WriteFileUser(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func WriteFileRoot(path string, data []byte) error {
	cmd := exec.Command("sudo", "-k", "tee", path)
	cmd.Stdin = bytes.NewReader(data)
	return cmd.Run()
}

func BackupFile(path string) (string, error) {
	backupPath := fmt.Sprintf("%s.bak.%s", path, time.Now().Format("20060102-150405"))
	data, err := os.ReadFile(path)
	if err != nil { return "", err }
	return backupPath, os.WriteFile(backupPath, data, 0o644)
}

func IsPermissionError(err error) bool {
	if err == nil { return false }
	if os.IsPermission(err) { return true }
	if errors.Is(err, syscall.EACCES) { return true }
	return false
}
