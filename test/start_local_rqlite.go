package test

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/jtarchie/forum/db"
	"github.com/onsi/gomega/gexec"
	"github.com/phayes/freeport"
)

func StartLocalRQLite() (*gexec.Session, *db.Client, error) {
	rqlitePath, err := os.MkdirTemp("", "")
	if err != nil {
		return nil, nil, fmt.Errorf("could not create temp directory: %w", err)
	}

	ports, err := freeport.GetFreePorts(2)
	if err != nil {
		return nil, nil, fmt.Errorf("could not find ports: %w", err)
	}

	rqliteCommand := exec.Command("rqlited",
		"-node-id", fmt.Sprintf("namespace-%d", rand.Int()),
		"-http-addr", fmt.Sprintf("localhost:%d", ports[0]),
		"-raft-addr", fmt.Sprintf("localhost:%d", ports[1]),
		rqlitePath,
	)

	rqliteSession, err := gexec.Start(rqliteCommand, nil, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("could not start execution: %w", err)
	}

	select {
	case <-rqliteSession.Err.Detect("entering leader state"):
	case <-time.After(10 * time.Second):
		return nil, nil, fmt.Errorf("could not start rqlite")
	}

	client, err := db.NewClient(fmt.Sprintf("http://localhost:%d", ports[0]))
	if err != nil {
		return nil, nil, fmt.Errorf("could not create client: %w", err)
	}

	return rqliteSession, client, nil
}
