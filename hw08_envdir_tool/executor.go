package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env []string) int {
	cm := cmd[0]
	args := cmd[1:]
	ex := exec.Command(cm, args...)

	ex.Stderr = os.Stderr
	ex.Stdin = os.Stdin
	ex.Stdout = os.Stdout
	ex.Env = env

	ex.Run()

	return ex.ProcessState.ExitCode()
}
