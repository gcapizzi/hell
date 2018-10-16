package main_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	cpuSetName := "test"

	AfterEach(func() {
		os.RemoveAll(cpusetPath(cpuSetName))
	})

	Context("when the cpuset already exists", func() {
		BeforeEach(func() {
			Expect(os.Mkdir(cpusetPath(cpuSetName), 0755)).To(Succeed())
			cmd := exec.Command(pinCpuPath, "-cpuset", cpuSetName, "-cpus", "0")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

			Eventually(session).Should(gexec.Exit(0))
			Expect(err).NotTo(HaveOccurred())
		})

		It("adds a CPU to a CPUset", func() {
			Expect(ioutil.ReadFile(cpusetPath(cpuSetName, "cpuset.cpus"))).To(BeEquivalentTo("0\n"))
		})

		It("copies memory from the parent cgroup", func() {
			parentMemoryFileContent, err := ioutil.ReadFile(cpusetPath("cpuset.mems"))
			Expect(err).NotTo(HaveOccurred())
			Expect(ioutil.ReadFile(cpusetPath(cpuSetName, "cpuset.mems"))).To(BeEquivalentTo(parentMemoryFileContent))
		})
	})

	Context("when the cpuset does not already exist", func() {
		It("creates a new cgroup", func() {
			cmd := exec.Command(pinCpuPath, "-cpuset", cpuSetName, "-cpus", "0")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

			Eventually(session).Should(gexec.Exit(0))
			Expect(err).NotTo(HaveOccurred())
			Expect(cpusetPath(cpuSetName)).To(BeADirectory())
		})
	})
})

func cpusetPath(parts ...string) string {
	return filepath.Join(append([]string{"/sys/fs/cgroup/cpuset"}, parts...)...)
}
