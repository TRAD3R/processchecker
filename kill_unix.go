//go:build !windows

package processchecker

import (
	"fmt"
	"path"
	"syscall"
)

func checkProcess(pidFile string) (string, error) {
	filepath := path.Join("/tmp", pidFile)
	// Проверка, активен ли процесс с этим PID
	if _, err := os.Stat(filepath); err != nil {
		return "", fmt.Errorf("cannot find process \"%s\"", filepath)
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

	err = syscall.Kill(pid, 0)
	if err == nil {
		return "", fmt.Errorf("service already running with PID %d", pid)
	}

	if err != syscall.ESRCH {
		return "", fmt.Errorf("failed to check PID %d: %w", pid, err)
	}

	return filepath, nil
}
