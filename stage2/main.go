package main

import (
	"fmt"
	"os"
	"os/exec"
)

// docker 			run image	<cmd> <args>
// go run main.go 	run 		<cmd> <args>

func main() {
	switch os.Args[1] {
	case "run":
		run()

	default:
		panic("bad command")
	}
}

func run() {
	fmt.Printf("Running %v \n", os.Args[2:])
	command := os.Args[2]
	args := os.Args[3:]
	cmd := exec.Command(command, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
