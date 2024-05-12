package main

import (
	"fmt"
	"os"
)

func main() {
	// Place your code here.
	path := os.Args[1]
	command := os.Args[2:]

	env, err := ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}

	os.Exit(RunCmd(command, env))
}
