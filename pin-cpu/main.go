package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	var cpusetName string
	var coreIndexes string

	flag.StringVar(&cpusetName, "cpuset", "", "the cpuset name")
	flag.StringVar(&coreIndexes, "cpus", "", "the core indexes")

	flag.Parse()

	os.Mkdir(cpusetPath(cpusetName), 0755)

	err := copyFile(cpusetPath("cpuset.mems"), cpusetPath(cpusetName, "cpuset.mems"))
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(cpusetPath(cpusetName, "cpuset.cpus"), []byte(coreIndexes), 0755)
	if err != nil {
		panic(err)
	}
}

func copyFile(src, dest string) error {
	parentMemsContent, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, []byte(parentMemsContent), 0755)
	if err != nil {
		return err
	}
	return nil
}

func cpusetPath(parts ...string) string {
	return filepath.Join(append([]string{"/sys/fs/cgroup/cpuset"}, parts...)...)
}
