package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

func main() {
	var lowerDirName string
	var separateFsDir string

	flag.StringVar(&lowerDirName, "rootfs", "", "the root fs")
	flag.StringVar(&separateFsDir, "graph", "", "the graph")

	flag.Parse()

	separateFsImage, err := createImageFile(1024 * 1024 * 1024)
	if err != nil {
		panic(err)
	}

	err = mountFsImage(separateFsImage, separateFsDir)
	if err != nil {
		panic(err)
	}

	workDirName, err := makeDir(separateFsDir)
	if err != nil {
		panic(err)
	}

	upperDirName, err := makeDir(separateFsDir)
	if err != nil {
		panic(err)
	}

	mergedDirName, err := setupOverlay(lowerDirName, upperDirName, workDirName)
	if err != nil {
		panic(err)
	}

	fmt.Println(mergedDirName)
}

func createImageFile(size int64) (*os.File, error) {
	tmpFile, err := ioutil.TempFile("", "image_*.img")
	if err != nil {
		return nil, err
	}

	err = unix.Fallocate(int(tmpFile.Fd()), 0, 0, size)
	if err != nil {
		return nil, err
	}

	err = run("mkfs.ext4", tmpFile.Name())
	if err != nil {
		return nil, err
	}

	return tmpFile, nil
}

func mountFsImage(image *os.File, dir string) error {
	err := run("mount", "-o", "loop", image.Name(), dir)
	if err != nil {
		return err
	}

	return nil
}

func setupOverlay(lowerDirName, upperDirName, workDirName string) (string, error) {
	mergedDirName, err := makeDir("")
	if err != nil {
		return "", err
	}

	options := fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDirName, upperDirName, workDirName)
	err = run("mount", "-t", "overlay", "overlay", "-o", options, mergedDirName)
	if err != nil {
		return "", err
	}

	return mergedDirName, nil
}

func makeDir(parentDirName string) (string, error) {
	dirName, err := ioutil.TempDir(parentDirName, "")
	if err != nil {
		return "", err
	}
	return dirName, nil

}

func run(cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	return command.Run()
}
