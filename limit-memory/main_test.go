package main_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	cgroupName := "test"

	AfterEach(func() {
		Expect(os.RemoveAll(memoryCgroupPath(cgroupName))).To(Succeed())
	})

	Context("when the cgroup already exists", func() {
		var sleepSession *gexec.Session
		var err error
		BeforeEach(func() {
			Expect(os.Mkdir(memoryCgroupPath(cgroupName), 0755)).To(Succeed())
			sleepSession, err = gexec.Start(exec.Command("sleep", "10"), GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Expect(ioutil.WriteFile(memoryCgroupPath(cgroupName, "tasks"), []byte(strconv.Itoa(sleepSession.Command.Process.Pid)), 0755)).To(Succeed())
			cmd := exec.Command(limitMemoryPath, "-cgroup", cgroupName, "-max", "42M")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

			Eventually(session).Should(gexec.Exit(0))
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			sleepSession.Kill().Wait(1)
		})

		It("overwrites the memory limit", func() {
			Expect(ioutil.ReadFile(memoryCgroupPath(cgroupName, "memory.limit_in_bytes"))).To(BeEquivalentTo("44040192\n"))
		})
	})

})

func memoryCgroupPath(parts ...string) string {
	return filepath.Join(append([]string{"/sys/fs/cgroup/memory"}, parts...)...)
}
