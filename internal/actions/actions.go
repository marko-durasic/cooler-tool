package actions

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
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

	// Use WaitGroup to parallelize governor setting for all cores
	var wg sync.WaitGroup
	wg.Add(numCPUs)

	for i := 0; i < numCPUs; i++ {
		go func(coreID int) {
			defer wg.Done()
			cmd := exec.Command("sudo", "cpufreq-set", "-c", strconv.Itoa(coreID), "-g", governor)
			_ = cmd.Run() // Errors are acceptable per-core
		}(i)
	}

	wg.Wait()
	fmt.Println("CPU governor update command sent.")
	return nil
}
