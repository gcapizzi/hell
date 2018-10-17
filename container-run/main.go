package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	if os.Args[1] == "child" {
		child()
	} else {
		parent()
	}
}

func parent() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[1:]...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS,
	}

	err := cmd.Run()
	if err != nil {
		exitWithError(err)
	}
}

func child() {
	var cpusetName string
	var rootfsPath string

	flagSet := flag.NewFlagSet("child", flag.ExitOnError)
	flagSet.StringVar(&cpusetName, "cpuset", "", "the cpuset name")
	flagSet.StringVar(&rootfsPath, "rootfs", "", "the rootfs path")
	flagSet.Parse(os.Args[2:])

	if rootfsPath == "" {
		panic("-rootfs is Required")
	}

	cmd := exec.Command(flagSet.Args()[0], flagSet.Args()[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	pid := os.Getpid()
	ioutil.WriteFile(filepath.Join("/sys/fs/cgroup/cpuset", cpusetName, "tasks"), []byte(strconv.Itoa(pid)), 0755)

	oldRootfsPath := filepath.Join(rootfsPath, "oldrootfs")
	must(syscall.Mount(rootfsPath, rootfsPath, "", syscall.MS_BIND, ""))
	must(os.MkdirAll(oldRootfsPath, 0700))
	must(syscall.PivotRoot(rootfsPath, oldRootfsPath))
	must(os.Chdir("/"))

	err := cmd.Run()
	if err != nil {
		exitWithError(err)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func exitWithError(err error) {
	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			os.Exit(status.ExitStatus())
		}
	} else {
		panic(err)
	}
}
