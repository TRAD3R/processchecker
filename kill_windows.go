//go:build windows

package processchecker

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"syscall"

	"golang.org/x/sys/windows"
)

const (
	ERROR_INVALID_PARAMETER = 87
)

func checkProcess(pidFile string) (string, error) {
	filepath := path.Join("%TMP%", pidFile)
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

	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		if errno, ok := err.(syscall.Errno); ok {
			// Процесс с таким PID не существует
			if errno == ERROR_INVALID_PARAMETER {
				return filepath, nil
			}
		}
		return "", fmt.Errorf("failed to check PID %d: %w", pid, err)
	}
	defer func() {
		if err := windows.CloseHandle(handle); err != nil {
			log.Printf("Error closing process handle: %v\n", err)
		}
	}()

	// Если OpenProcess успешно, процесс запущен
	return "", fmt.Errorf("service already running with PID %d", pid)
}
