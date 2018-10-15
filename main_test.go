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
})
