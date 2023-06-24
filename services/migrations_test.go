package services_test

import (
	"github.com/jtarchie/forum/db"
	"github.com/jtarchie/forum/services"
	"github.com/jtarchie/forum/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"go.uber.org/zap"
)

var _ = Describe("Migrations", func() {
	var (
		client        *db.Client
		rqliteSession *gexec.Session
	)

	BeforeEach(func() {
		var err error
		rqliteSession, client, err = test.StartLocalRQLite()
		Expect(err).NotTo(HaveOccurred())
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

			rerunResults, err := client.Query("SELECT * FROM migrations")
			Expect(err).NotTo(HaveOccurred())
			Expect(rerunResults.NumRows()).To(Equal(results.NumRows()))
		})
	})
})
