package cmd_test

import (
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/jtarchie/forum/cmd"
	"github.com/jtarchie/forum/db"
	"github.com/jtarchie/forum/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/phayes/freeport"
)

var _ = Describe("Server", func() {
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

	It("can accepts HTTP requests", func() {
		port := freeport.GetPort()
		command := cmd.ServerCmd{
			Port:     port,
			DBServer: client.URL(),
		}

		go func() {
			defer GinkgoRecover()
			err := command.Run()
			Expect(err).NotTo(HaveOccurred())
		}()

		response, err := req.C().R().SetRetryCount(3).Get(fmt.Sprintf("http://localhost:%d/", port))
		Expect(err).NotTo(HaveOccurred())

		Expect(response.ToString()).To(ContainSubstring("List Forums"))
	})
})
