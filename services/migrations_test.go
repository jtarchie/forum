package services_test

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/jtarchie/forum/db"
	"github.com/jtarchie/forum/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/phayes/freeport"
	"go.uber.org/zap"
)

var _ = Describe("Migrations", func() {
	var (
		client        *db.Client
		rqliteSession *gexec.Session
	)

	BeforeEach(func() {
		rqlitePath, err := os.MkdirTemp("", "")
		Expect(err).ShouldNot(HaveOccurred())

		port, err := freeport.GetFreePort()
		Expect(err).ShouldNot(HaveOccurred())

		rqliteCommand := exec.Command("rqlited",
			"-node-id", fmt.Sprintf("namespace-%d", GinkgoParallelProcess()),
			"-http-addr", fmt.Sprintf("localhost:%d", port),
			rqlitePath,
		)

		rqliteSession, err = gexec.Start(rqliteCommand, GinkgoWriter, GinkgoWriter)
		Expect(err).ShouldNot(HaveOccurred())
		Eventually(rqliteSession.Err, "10s").Should(gbytes.Say("entering leader state"))

		client, err = db.NewClient(fmt.Sprintf("http://localhost:%d", port))
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		rqliteSession.Kill()
	})

	When("migrations are run", func() {
		It("saves the information in the migrations table", func() {
			logger, err := zap.NewDevelopment()
			Expect(err).NotTo(HaveOccurred())

			err = services.Migration(client, logger)
			Expect(err).NotTo(HaveOccurred())

			results, err := client.Query("SELECT * FROM migrations")
			Expect(err).NotTo(HaveOccurred())
			Expect(results.NumRows()).To(BeNumerically(">", 0))

			By("ensuring they are idempotent")
			err = services.Migration(client, logger)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
