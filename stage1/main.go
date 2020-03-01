package main

import (
	"fmt"
	"os"
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
	fmt.Printf("Running %v\n", os.Args[2:])
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
