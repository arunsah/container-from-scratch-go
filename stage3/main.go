package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// docker 			run image	<cmd> <args>
// go run main.go 	run 		<cmd> <args>

func main() {
	fmt.Printf("args main:: %v\n", os.Args[0:])
	switch os.Args[1] {
	case "run":
		run()

	case "child":
		child()

	default:
		panic("bad command")
	}
}

func run() {
	// the will not run arbitary commands insted it will run itself child().
	// here we will setup the namesapce

	fmt.Printf("Running main:: %v\n", os.Args[2:])

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// stage 3
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS, // Unix time sharing system namespace
		// is the hostname. It will help us to provide out own hostname
		// inside container
	}

	cmd.Run()
}

func child() {
	// stage3b
	// this will help us to setup the hostname of container
	// run() will set the environment
	// child() will be actually be running in the container environemnt set
	//	setup by run()

	fmt.Printf("Running child:: %v \n", os.Args[2:])
	// by this time we are already in new namesapce
	// so we can set the container hostname
	syscall.Sethostname([]byte("container"))

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
