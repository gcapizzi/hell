package main_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Main", func() {
	var hellPath string

	BeforeEach(func() {
		var err error
		hellPath, err = gexec.Build("github.com/gcapizzi/hell")
		Expect(err).NotTo(HaveOccurred())
	})

	It("runs an arbitrary command", func() {
		session, err := gexec.Start(exec.Command("sudo", hellPath, "echo", "bye, world!"), GinkgoWriter, GinkgoWriter)

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
			cmd := exec.Command("sudo", hellPath, "bash", "-c", "hostname foo; hostname")

			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)

			Eventually(session).Should(gexec.Exit())
			Expect(err).NotTo(HaveOccurred())
			Expect(session.Out).To(gbytes.Say("foo"))

			command := exec.Command("hostname")
			session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit())

			Expect(session.Out).To(gbytes.Say(systemHostname))
		})
	})
})
