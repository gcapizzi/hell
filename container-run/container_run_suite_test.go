package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestContainerRun(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ContainerRun Suite")
}

var hellPath string

var _ = BeforeSuite(func() {
	var err error
	hellPath, err = gexec.Build("github.com/gcapizzi/hell/container-run")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
