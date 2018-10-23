package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestLimitMemory(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LimitMemory Suite")
}

var limitMemoryPath string

var _ = BeforeSuite(func() {
	var err error
	limitMemoryPath, err = gexec.Build("github.com/gcapizzi/hell/limit-memory")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
