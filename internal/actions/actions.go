package actions

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"runtime"
)

func KillProcess(pid string) error {
	pidInt, err := strconv.Atoi(pid)
	if err != nil {
		return fmt.Errorf("invalid PID")
	}
	process, err := os.FindProcess(pidInt)
	if err != nil {
		return err
	}
	return process.Kill()
}

func SetCpuGovernor(governor string) error {
	numCPUs := runtime.NumCPU()
	fmt.Printf("Setting CPU governor to '%s' for %d cores...\n", governor, numCPUs)
	for i := 0; i < numCPUs; i++ {
		cmd := exec.Command("sudo", "cpufreq-set", "-c", strconv.Itoa(i), "-g", governor)
		err := cmd.Run()
		if err != nil {
			// Log or handle error per core if needed, for now we just continue
		}
	}
	fmt.Println("CPU governor update command sent.")
	return nil
}
