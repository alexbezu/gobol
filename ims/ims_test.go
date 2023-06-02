//go:build immudb

package ims

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func provision_immudb(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "alexbezuglyi/mqimmudb:latest",
		Networks:     []string{"mynet"},
		Hostname:     "db",
		ExposedPorts: []string{"3322/tcp"},
		// WaitingFor:   wait.ForListeningPort("3322/tcp"),
		WaitingFor: wait.ForLog("Admin user 'immudb' successfully created"),
	}
	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func provision_transaction(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context: "testcontainer", // this is the folder inside the ims package
		},
		Networks:   []string{"mynet"},
		WaitingFor: wait.ForLog("Connected!"),
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func Test3270(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	ims_cont, err := provision_immudb(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	tran, err := provision_transaction(ctx)
	if err != nil {
		t.Error(err)
		return
	}
	ep, err := ims_cont.Endpoint(ctx, "")
	if err != nil {
		t.Error(err)
		return
	}

	Connect2db(ep)
	defer closedb()

	screen := NewTN3270screen()
	os.Setenv("TN3270DIR", "testcontainer")
	screen.Lterm = "testing1"
	screen.readFormat("SCREEN")
	(*screen.DFLD)[screen.POS(4, 25)].value.Set("Ping")
	screen.MFLDsend(0x7d)
	screen.MFLDrecieve()
	if strings.Trim(screen.MFLDout.I["RESULT"].String(), string(byte(0x00))) != "Pong" {
		t.Error("Bad response from transaction")
	}

	defer func() {
		if err := ims_cont.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()
	defer func() {
		if err := tran.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()
}
