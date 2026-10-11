//go:build !windows

package terminal

import (
	"os"
	"syscall"
)

type shellSignal = syscall.Signal

const (
	sigHangup = syscall.SIGHUP
	sigKill   = syscall.SIGKILL
)

// go-pty starts the shell with Setsid, so its pid is also its process group id.
// An interactive shell forwards SIGHUP to its jobs, which live in their own groups.
func signalProcessGroup(p *os.Process, sig shellSignal) {
	if err := syscall.Kill(-p.Pid, sig); err != nil {
		_ = p.Signal(sig)
	}
}
