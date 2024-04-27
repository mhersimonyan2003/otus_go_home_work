package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Place your code here.
	// get command line arguments
	path := os.Args[1]
	command := os.Args[2:]

	env, err := ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}
	RunCmd(command, env)

	cmd := exec.Command("/bin/bash", "./testdata/echo.sh", "arg1=1", "arg2=2")
	cmd.Env = os.Environ()
	if err != nil {
		fmt.Println("Error:", err)
	}
}
