package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestPinCpu(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PinCpu Suite")
}

var pinCpuPath string

var _ = BeforeSuite(func() {
	var err error
	pinCpuPath, err = gexec.Build("github.com/gcapizzi/hell/pin-cpu")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
