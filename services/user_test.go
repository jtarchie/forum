package services_test

import (
	"github.com/jtarchie/forum/db"
	"github.com/jtarchie/forum/services"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
)

var _ = Describe("User", func() {
	var (
		client db.Client
	)

	BeforeEach(func() {
		var err error

		client, err = db.NewClient("sqlite://")
		Expect(err).NotTo(HaveOccurred())

		logger, err := zap.NewDevelopment()
		Expect(err).NotTo(HaveOccurred())

		err = services.Migration(client, logger)
		Expect(err).NotTo(HaveOccurred())
	})

	userCount := func(client db.Client) int {
		rows, err := client.Query("SELECT COUNT(*) FROM users;")
		Expect(err).NotTo(HaveOccurred())
		defer rows.Close()

		var count int
		Expect(rows.Next()).To(BeTrue())

		err = rows.Scan(&count)
		Expect(err).NotTo(HaveOccurred())

		return count
	}

	When("user does not exist", func() {
		It("creates one", func() {
			Expect(userCount(client)).To(Equal(0))

			err := services.UpsertUser(client, services.User{
				Email:    "bob@example.com",
				Provider: "github",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(userCount(client)).To(Equal(1))
		})
	})

	When("the user already exists", func() {
		It("does nothing", func() {
			Expect(userCount(client)).To(Equal(0))

			err := services.UpsertUser(client, services.User{
				Email:    "bob@example.com",
				Provider: "github",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(userCount(client)).To(Equal(1))

			err = services.UpsertUser(client, services.User{
				Email:    "bob@example.com",
				Provider: "github",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(userCount(client)).To(Equal(1))

			By("setting a random username")

			rows, err := client.Query("SELECT username FROM users;")
			Expect(err).NotTo(HaveOccurred())
			defer rows.Close()

			var username string
			Expect(rows.Next()).To(BeTrue())

			err = rows.Scan(&username)
			Expect(err).NotTo(HaveOccurred())
			Expect(username).To(MatchRegexp(`\w+-\w+-\d{4}`))
		})
	})
})
