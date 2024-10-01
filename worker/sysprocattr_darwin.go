//go:build darwin

package worker

import "syscall"

const bash = "bash"

func getDaemonSysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		Setsid: true,
	}
}
