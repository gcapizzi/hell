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
	var cpusetName string

	flag.StringVar(&cpusetName, "cpuset", "", "the cpuset name")
	flag.Parse()

	cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}

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
