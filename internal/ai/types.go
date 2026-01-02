package ai

type ActionType string

const (
    InstallPackage ActionType = "install_package"
    RemovePackage  ActionType = "remove_package"
    ReadFile       ActionType = "read_file"
    EditFile       ActionType = "edit_file"
    CreateFile     ActionType = "create_file"
    RunCommand     ActionType = "run_command"
)

type Action struct {
    Type ActionType `json:"type"`
    Package string `json:"package,omitempty"`
    Path string `json:"path,omitempty"`
    Language string `json:"language,omitempty"`
    Content string `json:"content,omitempty"`
    Diff string `json:"diff,omitempty"`
    Command []string `json:"command,omitempty"`
    NeedsRoot bool `json:"needs_root,omitempty"`
}

type Plan struct {
    Actions []Action `json:"actions"`
    Explanation string `json:"explanation"`
}
