package main

import (
	"context"
	"log"
	"math/rand/v2"
	"os"
	"runtime/debug"
	"strings"

	"github.com/neatflowcv/cepher/pkg/cephcli"
	"github.com/neatflowcv/cepher/pkg/cephrest"
	"github.com/neatflowcv/cepher/pkg/cephsetup"
)

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	return info.Main.Version
}

func Main2() {
	log.Println("version", version())

	apiURL := os.Getenv("CEPH_API_URL")
	username := os.Getenv("CEPH_USERNAME")
	password := os.Getenv("CEPH_PASSWORD")

	if apiURL == "" || username == "" || password == "" {
		log.Fatalf("CEPH_API_URL, CEPH_USERNAME, and CEPH_PASSWORD must be set")
	}

	client := cephrest.NewClient(apiURL, username, password)

	authResponse, err := client.Auth(context.Background())
	if err != nil {
		log.Fatalf("failed to authenticate with ceph: %v", err)
	}

	err = client.GetHealthFull(context.Background(), authResponse.Token)
	if err != nil {
		log.Fatalf("failed to get cluster from ceph: %v", err)
	}

	err = client.Logout(context.Background(), authResponse.Token)
	if err != nil {
		log.Fatalf("failed to logout from ceph: %v", err)
	}

	log.Println("logged out from ceph")
}

func main() {
	log.Println("version", version())

	hosts := os.Getenv("CEPH_MON_HOSTS")
	if hosts == "" {
		log.Panicf("CEPH_MON_HOSTS must be set")
	}

	keyring := os.Getenv("CEPH_KEYRING")
	if keyring == "" {
		log.Panicf("CEPH_KEYRING must be set")
	}

	monHosts := strings.Split(hosts, ",")
	rand.Shuffle(len(monHosts), func(i, j int) {
		monHosts[i], monHosts[j] = monHosts[j], monHosts[i]
	})

	tempDir, err := os.MkdirTemp("", "cepher")
	if err != nil {
		log.Panicf("failed to create temporary directory: %v", err)
	}

	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			log.Printf("failed to remove temporary directory: %v", err)
		}
	}()

	err = cephsetup.Setup(tempDir, monHosts, keyring)
	if err != nil {
		log.Panicf("failed to setup ceph: %v", err)
	}

	client := cephcli.NewClient(tempDir, "20.1.1")

	health, err := client.GetHealth(context.Background())
	if err != nil {
		log.Panicf("failed to get health detail from ceph: %v", err)
	}

	log.Println("health", health)
}
