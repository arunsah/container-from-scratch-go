package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

// docker 			run image	<cmd> <args>
// go run main.go 	run 		<cmd> <args>
// video 26 mins/43 mins
// chrtoot, proccess
// control group limits the resources
// CGroups: what you can use
//	 file system interface: memory/cpu/io/process number,...

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

	fmt.Printf("Running %v as %d\n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// stage 3
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,

		//CLONE_NEWUTS: Unix time sharing system namespace
		// is the hostname. It will help us to provide out own hostname
		// inside container

		// CLONE_NEWPID: process id namespace
		// CLONE_NEWNS: first namespae (for mount) created on linux

		// mounts are by default, shared by all (host) namespaces
		// after this our mounts within container will not be shared by any other namesapce
		Unshareflags: syscall.CLONE_NEWNS,
	}

	cmd.Run()
}

func child() {
	// stage3b
	// this will help us to setup the hostname of container
	// run() will set the environment
	// child() will be actually be running in the container environemnt set
	//	setup by run()

	fmt.Printf("Running child %v as %d\n", os.Args[2:], os.Getpid())

	//stage5
	cg()

	// by this time we are already in new namesapce
	// so we can set the container hostname
	syscall.Sethostname([]byte("container"))

	// stage 4
	//syscall.Chroot("/garlicbread/") // this will be root dir for out container
	//syscall.Chdir("/")              // cd else undefined behaviour

	if err := syscall.Chroot("/garlicbread/minibuntu"); err != nil {
		panic(err)
	}

	if err := os.Chdir("/"); err != nil {
		panic(err)
	}

	// stage4b: after this 'ps' will work
	syscall.Mount("proc", "proc", "proc", 0, "")

	command := os.Args[2]
	args := os.Args[3:]
	cmd := exec.Command(command, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//must(cmd.Run())
	cmd.Run()

	// stage4b
	syscall.Unmount("/proc", 0)
}

func cg() {

	cgroup := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroup, "pids")
	err := os.Mkdir(filepath.Join(pids, "garlic"), 0755)

	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	must(ioutil.WriteFile(filepath.Join(pids, "garlic/pids.max"), []byte("14"), 0700))
	// remove the new cgroup in place after the container exits
	must(ioutil.WriteFile(filepath.Join(pids, "garlic/notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "garlic/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))

}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
