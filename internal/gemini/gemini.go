package gemini

import (
	"cooler/internal/diagnostics"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// Cache the project directory to avoid repeated os.Executable() calls
	projectDirOnce sync.Once
	cachedProjectDir string
	projectDirErr    error
)

func getProjectDir() (string, error) {
	projectDirOnce.Do(func() {
		exePath, err := os.Executable()
		if err != nil {
			projectDirErr = fmt.Errorf("error getting executable path: %w", err)
			return
		}
		cachedProjectDir = filepath.Dir(exePath)
	})
	return cachedProjectDir, projectDirErr
}

// AskGemini executes the gemini cli and returns the analysis as a string.
func AskGemini(maxTemp float64, topProcess *diagnostics.ProcessInfo) (string, error) {
	projectDir, err := getProjectDir()
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(
		`You are an expert system analyst. I am a script providing you with data about a user's computer that is overheating. `+
			`Please analyze the following information and provide a brief, user-friendly diagnosis and suggestion. `+
			`Do not ask questions, provide a direct analysis. Keep the response to a few sentences.\n\n`+
			`Data:\n- Max CPU Temperature: %.1f°C\n`+
			`- Top CPU Process Name: %s\n`+
			`- Top Process PID: %s\n`+
			`- Top Process CPU %%: %s\n\n`+
			`Analysis:`, 
		maxTemp, topProcess.Name, topProcess.PID, topProcess.CPU,
	)

	cmd := exec.Command("gemini", "--prompt", prompt)
	cmd.Dir = projectDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(out), "command not found") {
			return "", fmt.Errorf("the 'gemini' command was not found")
		}
		return "", fmt.Errorf("error calling Gemini: %w, output: %s", err, string(out))
	}

	return string(out), nil
}