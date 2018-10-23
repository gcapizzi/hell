package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	var cgroupName string
	var maxMemory string

	flag.StringVar(&cgroupName, "cgroup", "", "the cgroup name")
	flag.StringVar(&maxMemory, "max", "", "the max memory in bytes")

	flag.Parse()

	os.Mkdir(cgroupPath(cgroupName), 0755)
	err := ioutil.WriteFile(cgroupPath(cgroupName, "memory.limit_in_bytes"), []byte(maxMemory), 0755)
	if err != nil {
		panic(err)
	}
}

func cgroupPath(parts ...string) string {
	return filepath.Join(append([]string{"/sys/fs/cgroup/memory"}, parts...)...)
}
