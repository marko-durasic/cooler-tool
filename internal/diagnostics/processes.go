package diagnostics

import (
	"fmt"
	"os/exec"
	"strings"
)

type ProcessInfo struct {
	PID  string
	CPU  string
	Name string
}

func GetTopProcessInfo() (*ProcessInfo, error) {
	cmd := exec.Command("bash", "-c", "ps -eo pcpu,pid,comm --sort=-pcpu | head -n 2 | tail -n 1")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	output := strings.TrimSpace(string(out))
	parts := strings.Fields(output)
	if len(parts) < 3 {
		return nil, fmt.Errorf("unexpected output from ps command")
	}

	return &ProcessInfo{
		CPU:  parts[0],
		PID:  parts[1],
		Name: strings.Join(parts[2:], " "),
	}, nil
}
