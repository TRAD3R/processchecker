//go:build !windows

package processchecker

import (
	"fmt"
	"path"
	"syscall"

	"errors"
	"os"
	"strconv"
)

func checkProcess(pidFile string) (string, error) {
	filepath := path.Join("/tmp", pidFile)
	if _, err := os.Stat(filepath); err != nil {
		return filepath, nil
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read PID file: %w", err)
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		if errR := os.Remove(filepath); errR != nil {
			err = errors.Join(err, fmt.Errorf("failed to remove PID file %s: %w", pidFile, errR))
		}

		return "", fmt.Errorf("iInvalid PID in PID file: %w", err)
	}

	// Проверка, активен ли процесс с этим PID
	err = syscall.Kill(pid, 0)
	if err == nil {
		return "", fmt.Errorf("service already running with PID %d", pid)
	}

	if err != syscall.ESRCH {
		return "", fmt.Errorf("failed to check PID %d: %w", pid, err)
	}

	return filepath, nil
}
