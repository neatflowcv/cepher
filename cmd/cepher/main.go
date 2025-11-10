package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"

	"github.com/neatflowcv/cepher/pkg/cephrest"
)

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown"
	}

	return info.Main.Version
}

func main() {
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
