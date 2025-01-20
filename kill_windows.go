//go:build windows

package processchecker

import (
	"fmt"
	"log"
	"syscall"

	"golang.org/x/sys/windows"
)

const (
	ERROR_INVALID_PARAMETER = 87
)

func checkProcess(pid int) error {
	handle, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(pid))
	if err != nil {
		if errno, ok := err.(syscall.Errno); ok {
			// Процесс с таким PID не существует
			if errno == ERROR_INVALID_PARAMETER {
				return nil
			}
		}
		return fmt.Errorf("failed to check PID %d: %w", pid, err)
	}
	defer func() {
		if err := windows.CloseHandle(handle); err != nil {
			log.Printf("Error closing process handle: %v\n", err)
		}
	}()

	// Если OpenProcess успешно, процесс запущен
	return fmt.Errorf("service already running with PID %d", pid)
}
