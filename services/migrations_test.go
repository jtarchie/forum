package services_test

import (
	"github.com/jtarchie/forum/db"
	"github.com/jtarchie/forum/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("Migrations", func() {
	var (
		client db.Client
	)

	BeforeEach(func() {
		var err error

		client, err = db.NewClient("sqlite://")
		Expect(err).NotTo(HaveOccurred())
	})

	When("migrations are run", func() {
		It("saves the information in the migrations table", func() {
			logger, err := zap.NewDevelopment()
			Expect(err).NotTo(HaveOccurred())

			err = services.Migration(client, logger)
			Expect(err).NotTo(HaveOccurred())

			results, err := client.Query("SELECT * FROM migrations")
			Expect(err).NotTo(HaveOccurred())

			count := 0
			for results.Next() {
				count++
			}
			Expect(count).To(BeNumerically(">", 0))

			By("ensuring they are idempotent")
			err = services.Migration(client, logger)
			Expect(err).NotTo(HaveOccurred())

			rerunResults, err := client.Query("SELECT * FROM migrations")
			Expect(err).NotTo(HaveOccurred())

			postCount := 0
			for rerunResults.Next() {
				postCount++
			}
			Expect(postCount).To(Equal(count))
		})
	})
})
