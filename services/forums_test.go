package services_test

import (
	"github.com/jtarchie/forum/db"
	"github.com/jtarchie/forum/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"go.uber.org/zap"
)

var _ = Describe("Forums", func() {
	var (
		client        *db.Client
		rqliteSession *gexec.Session
	)

	BeforeEach(func() {
		rqliteSession, client = startRqlite()
	})

	AfterEach(func() {
		rqliteSession.Kill()
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
