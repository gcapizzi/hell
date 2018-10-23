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

	if cgroupName != "" {
		must(setupCGroup(cgroupName))
	}

	must(setupRootFS(rootfsPath))

	err := cmd.Start()
	if err != nil {
		exitWithError(err)
	}

	if cgroupName != "" {
		must(moveParentPidToRootCgroups())
	}

	err = cmd.Wait()
	if err != nil {
		exitWithError(err)
	}
}

func setupCGroup(name string) error {
	pid := os.Getpid()

	err := setupCpuset(name, pid)
	if err != nil {
		return err
	}

	err = setupMemory(name, pid)
	if err != nil {
		return err
	}

	return nil
}

func setupCpuset(name string, pid int) error {
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

	return ioutil.WriteFile(cgroupPath("cpuset", name, "tasks"), []byte(strconv.Itoa(pid)), 0755)
}

func setupMemory(name string, pid int) error {
	if _, err := os.Stat(cgroupPath("memory", name)); os.IsNotExist(err) {
		os.Mkdir(cgroupPath("memory", name), 0755)
	}

	return ioutil.WriteFile(cgroupPath("memory", name, "tasks"), []byte(strconv.Itoa(pid)), 0755)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func setupRootFS(path string) error {
	oldRootfsPath := filepath.Join(path, "oldrootfs")

	err := syscall.Mount(path, path, "", syscall.MS_BIND, "")
	if err != nil {
		return err
	}

	err = os.MkdirAll(oldRootfsPath, 0700)
	if err != nil {
		return err
	}

	err = syscall.PivotRoot(path, oldRootfsPath)
	if err != nil {
		return err
	}

	return os.Chdir("/")
}

func moveParentPidToRootCgroups() error {
	pid := os.Getpid()

	err := movePidToRootCgroup(pid, "cpuset")
	if err != nil {
		return err
	}

	err = movePidToRootCgroup(pid, "memory")
	if err != nil {
		return err
	}

	return nil
}

func movePidToRootCgroup(pid int, cgroupType string) error {
	return ioutil.WriteFile(filepath.Join("/oldrootfs/sys/fs/cgroup", cgroupType, "tasks"), []byte(strconv.Itoa(pid)), 0755)
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
