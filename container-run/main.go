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
	var cgroupName string
	var rootfsPath string

	flagSet := flag.NewFlagSet("child", flag.ExitOnError)
	flagSet.StringVar(&cgroupName, "cgroup", "", "the cgroup name")
	flagSet.StringVar(&rootfsPath, "rootfs", "", "the rootfs path")
	flagSet.Parse(os.Args[2:])

	if rootfsPath == "" {
		panic("-rootfs is Required")
	}

	cmd := exec.Command(flagSet.Args()[0], flagSet.Args()[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	must(setupCGroup(cgroupName))
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

func setupCGroup(name string) error {
	if _, err := os.Stat(cgroupPath("cpuset", name)); os.IsNotExist(err) {
		os.Mkdir(cgroupPath("cpuset", name), 0755)

		err := copyFile(cgroupPath("cpuset", "cpuset.cpus"), cgroupPath("cpuset", name, "cpuset.cpus"))
		if err != nil {
			return err
		}

		err = copyFile(cgroupPath("cpuset", "cpuset.mems"), cgroupPath("cpuset", name, "cpuset.mems"))
		if err != nil {
			return err
		}
	}

	pid := os.Getpid()
	return ioutil.WriteFile(cgroupPath("cpuset", name, "tasks"), []byte(strconv.Itoa(pid)), 0755)
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

func cgroupPath(parts ...string) string {
	return filepath.Join(append([]string{"/sys/fs/cgroup"}, parts...)...)
}

func copyFile(src, dest string) error {
	srcContents, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, []byte(srcContents), 0755)
	if err != nil {
		return err
	}

	return nil
}
