//go:build windows

package worker

import "syscall"

const bash = "git-bash.exe"

func getDaemonSysProcAttr() *syscall.SysProcAttr {
	const DETACHED_PROCESS = 0x00000008

	return &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | DETACHED_PROCESS,
	}
}
