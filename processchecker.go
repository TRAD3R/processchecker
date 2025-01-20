package processchecker

import (
	"fmt"
	"os"
	"strconv"
)

func Run(pidFile string) error {
	filepath, err := checkProcess(pidFile)
	if err != nil {
		return err
	}

	// Запись текущего PID в PID файл
	pid := os.Getpid()
	err = os.WriteFile(filepath, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return fmt.Errorf("Failed to write %d in PID file: %w\n", pid, err)
	}

	return nil
}
