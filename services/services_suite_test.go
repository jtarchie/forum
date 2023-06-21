package services_test

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/jtarchie/forum/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/phayes/freeport"
)

func TestServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Services Suite")
}

func startRqlite() (*gexec.Session, *db.Client) {
	rqlitePath, err := os.MkdirTemp("", "")
	Expect(err).ShouldNot(HaveOccurred())

	ports, err := freeport.GetFreePorts(2)
	Expect(err).ShouldNot(HaveOccurred())

	rqliteCommand := exec.Command("rqlited",
		"-node-id", fmt.Sprintf("namespace-%d", GinkgoParallelProcess()),
		"-http-addr", fmt.Sprintf("localhost:%d", ports[0]),
		"-raft-addr", fmt.Sprintf("localhost:%d", ports[1]),
		rqlitePath,
	)

	rqliteSession, err := gexec.Start(rqliteCommand, GinkgoWriter, GinkgoWriter)
	Expect(err).ShouldNot(HaveOccurred())
	Eventually(rqliteSession.Err, "10s").Should(gbytes.Say("entering leader state"))

	client, err := db.NewClient(fmt.Sprintf("http://localhost:%d", ports[0]))
	Expect(err).ShouldNot(HaveOccurred())

	return rqliteSession, client
}
