package diagnostics

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

// Pre-compiled regex for better performance - avoids recompilation on every call
var coreTemperatureRegex = regexp.MustCompile(`Core \d+:\s+\+([\d\.]+)°C`)

func GetMaxCpuTemperature() (float64, error) {
	cmd := exec.Command("sensors")
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	matches := coreTemperatureRegex.FindAllStringSubmatch(string(out), -1)

	var maxTemp float64
	if len(matches) == 0 {
		return 0, fmt.Errorf("no core temperatures found")
	}

	for i, match := range matches {
		temp, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			continue
		}
		if i == 0 || temp > maxTemp {
			maxTemp = temp
		}
	}
	return maxTemp, nil
}
