package main

import (
	"os"
	"os/exec"
	"strings"
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

	execute.Stdin = os.Stdin
	execute.Stdout = os.Stdout
	execute.Stderr = os.Stderr

	execute.Env = os.Environ()

	for key, envItem := range env {
		if envItem.NeedRemove {
			execute.Env = removeEnvKeyExecute(execute.Env, key)

			if envItem.Value != "" {
				execute.Env = append(execute.Env, key+"="+envItem.Value)
			}
		} else {
			execute.Env = append(execute.Env, key+"="+envItem.Value)
		}
	}

	execute.Run()

	return execute.ProcessState.ExitCode()
}

func removeEnvKeyExecute(env []string, key string) []string {
	for i, envItem := range env {
		if strings.Split(envItem, "=")[0] == key {
			env = append(env[:i], env[i+1:]...)
			break
		}
	}

	return env
}
