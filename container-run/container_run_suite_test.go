package main_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestContainerRun(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ContainerRun Suite")
}

var containerRunPath string
var pinCpuPath string
var rootfsPath string

var _ = BeforeSuite(func() {
	var err error
	containerRunPath, err = gexec.Build("github.com/gcapizzi/hell/container-run")
	Expect(err).NotTo(HaveOccurred())
	pinCpuPath, err = gexec.Build("github.com/gcapizzi/hell/pin-cpu")
	Expect(err).NotTo(HaveOccurred())
	rootfsPath = os.Getenv("ROOTFS_PATH")
	Expect(rootfsPath).NotTo(BeEmpty(), "ROOTFS_PATH is required for the tests")
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
