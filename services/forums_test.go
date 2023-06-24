package services_test

import (
	"github.com/jtarchie/forum/db"
	"github.com/jtarchie/forum/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("Forums", func() {
	var (
		client db.Client
	)

	BeforeEach(func() {
		var err error

		client, err = db.NewClient("sqlite://")
		Expect(err).NotTo(HaveOccurred())
	})

	It("returns a list of forums", func() {
		logger, err := zap.NewDevelopment()
		Expect(err).NotTo(HaveOccurred())

		_, err = services.ListForums(client)
		Expect(err).To(HaveOccurred())

		err = services.Migration(client, logger)
		Expect(err).NotTo(HaveOccurred())

		forums, err := services.ListForums(client)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(forums)).To(BeNumerically(">", 0))
	})
})
