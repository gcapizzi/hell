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
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		} else {
			panic(err)
		}
	}
}

func child() {
	var cpusetName string

	flagSet := flag.NewFlagSet("child", flag.ExitOnError)
	flagSet.StringVar(&cpusetName, "cpuset", "", "the cpuset name")
	flagSet.Parse(os.Args[2:])

	cmd := exec.Command(flagSet.Args()[0], flagSet.Args()[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	pid := os.Getpid()
	ioutil.WriteFile(filepath.Join("/sys/fs/cgroup/cpuset", cpusetName, "tasks"), []byte(strconv.Itoa(pid)), 0755)

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		} else {
			panic(err)
		}
	}
}
