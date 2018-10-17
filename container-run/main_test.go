package main_test

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	It("runs an arbitrary command", func() {
		session, err := gexec.Start(exec.Command(containerRunPath, "-rootfs", "busybox", "echo", "bye, world!"), GinkgoWriter, GinkgoWriter)

		Eventually(session).Should(gexec.Exit())
		Expect(err).NotTo(HaveOccurred())
		Expect(session.Out).To(gbytes.Say("bye, world!"))
	})

	Context("hostname isolation", func() {
		var systemHostname string

		BeforeEach(func() {
			command := exec.Command("hostname")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit())
			systemHostname = string(session.Out.Contents())
		})

		It("runs the command in its own UPS namespace", func() {
			cmd := exec.Command(containerRunPath, "-rootfs", "busybox", "sh", "-c", "hostname foo; hostname")

			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

			Eventually(session).Should(gexec.Exit(0))
			Expect(err).NotTo(HaveOccurred())
			Expect(session.Out).To(gbytes.Say("foo"))

			command := exec.Command("hostname")
			session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit())

			Expect(session.Out).To(gbytes.Say(systemHostname))
		})
	})

	Context("process exit code forwarding", func() {
		It("exits with the same code as the internal process", func() {
			cmd := exec.Command(containerRunPath, "-rootfs", "busybox", "sh", "-c", "exit 42")

			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

			Eventually(session).Should(gexec.Exit(42))
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("cgroup", func() {
		It("runs the command in the specified cgroup", func() {
			cmd := exec.Command(pinCpuPath, "-cpuset", "test", "-cpus", "0")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Eventually(session).Should(gexec.Exit(0))

			cmd = exec.Command(containerRunPath, "-cgroup", "test", "-rootfs", "busybox", "sleep", "2")
			session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			defer session.Interrupt()

			Expect(err).NotTo(HaveOccurred())
			Eventually(func() ([]byte, error) {
				return ioutil.ReadFile("/sys/fs/cgroup/cpuset/test/tasks")
			}).ShouldNot(BeEmpty())
		})

		It("sets up the cgroup if it doesn't exist", func() {
			cgroupName := strconv.FormatInt(time.Now().UnixNano(), 16)

			cmd := exec.Command(containerRunPath, "-cgroup", cgroupName, "-rootfs", "busybox", "sleep", "2")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			defer session.Interrupt()

			Expect(err).NotTo(HaveOccurred())
			Eventually(func() ([]byte, error) {
				return ioutil.ReadFile(filepath.Join("/sys/fs/cgroup/cpuset", cgroupName, "tasks"))
			}).ShouldNot(BeEmpty())

			Eventually(func() ([]byte, error) {
				return ioutil.ReadFile(filepath.Join("/sys/fs/cgroup/memory", cgroupName, "tasks"))
			}).ShouldNot(BeEmpty())
		})
	})
})
