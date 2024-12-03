package processchecker

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"syscall"
)

type Logger struct {
	*slog.Logger
}

func Run(pidFile string) error {
	// Шаг 1: Проверка существования PID файла
	if _, err := os.Stat(pidFile); err == nil {
		data, err := os.ReadFile(pidFile)
		if err != nil {
			return fmt.Errorf("failed to read PID file: %w", err)
		}
		pid, err := strconv.Atoi(string(data))
		if err != nil {
			if errR := os.Remove(pidFile); errR != nil {
				err = errors.Join(err, fmt.Errorf("failed to remove PID file %s: %w", pidFile, errR))
			}

			return fmt.Errorf("iInvalid PID in PID file: %w", err)
		}

		// Проверка, активен ли процесс с этим PID
		err = syscall.Kill(pid, 0)
		if err == nil {
			return fmt.Errorf("service already running with PID %d", pid)
		} else if err != syscall.ESRCH {
			return fmt.Errorf("failed to check PID %d: %w", pid, err)
		}
	}

	// Шаг 2: Запись текущего PID в PID файл
	pid := os.Getpid()
	err := os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		return fmt.Errorf("Failed to write %d in PID file: %v\n", pid, err)
	}

	return nil
}
