package main

import (
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Place your code here.
	if len(cmd) == 0 {
		return -1
	}

	command := cmd[0]
	args := cmd[1:]

	// Create the command
	execute := exec.Command(command, args...)
	execute.Stdout = os.Stdout

	for key, envItem := range env {
		if envItem.NeedRemove {
			os.Unsetenv(key)
			if envItem.Value != "" {
				os.Setenv(key, envItem.Value)
			}
		} else {
			os.Setenv(key, envItem.Value)
		}
	}

	execute.Env = os.Environ()
	err := execute.Run()
	if err != nil {
		fmt.Println(err)
		return -1
	}

	return execute.ProcessState.ExitCode()
}
