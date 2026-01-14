package diagnostics

import (
	"bufio"
	"bytes"
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
	// Execute ps directly without shell overhead (avoids bash -c and pipe subprocesses)
	cmd := exec.Command("ps", "-eo", "pcpu,pid,comm", "--sort=-pcpu")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse output directly - skip header (line 0), get first data line (line 1)
	scanner := bufio.NewScanner(bytes.NewReader(out))
	lineNum := 0
	for scanner.Scan() {
		if lineNum == 1 { // Second line (first data line after header)
			line := strings.TrimSpace(scanner.Text())
			parts := strings.Fields(line)
			if len(parts) < 3 {
				return nil, fmt.Errorf("unexpected output from ps command")
			}
			return &ProcessInfo{
				CPU:  parts[0],
				PID:  parts[1],
				Name: strings.Join(parts[2:], " "),
			}, nil
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("no process data found")
}
