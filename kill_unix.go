//go:build !windows

package processchecker

import (
	"fmt"
	"syscall"
)

func checkProcess(pid int) error {
	err := syscall.Kill(pid, 0)
	if err == nil {
		return fmt.Errorf("service already running with PID %d", pid)
	}

	if err != syscall.ESRCH {
		return fmt.Errorf("failed to check PID %d: %w", pid, err)
	}

	return nil
}
