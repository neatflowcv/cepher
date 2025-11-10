package main

import (
	"context"
	"errors"
	"fmt"
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

type CephCLIConfig struct {
	Hosts   []string
	Keyring string
}

var (
	ErrRequiredEnvironment = errors.New("required environment is not set")
)

func LoadCephCLIConfig() (*CephCLIConfig, error) {
	hosts := os.Getenv("CEPH_MON_HOSTS")
	if hosts == "" {
		return nil, fmt.Errorf("%w: CEPH_MON_HOSTS", ErrRequiredEnvironment)
	}

	keyring := os.Getenv("CEPH_KEYRING")
	if keyring == "" {
		return nil, fmt.Errorf("%w: CEPH_KEYRING", ErrRequiredEnvironment)
	}

	return &CephCLIConfig{
		Hosts:   strings.Split(hosts, ","),
		Keyring: keyring,
	}, nil
}

func main() {
	log.Println("version", version())

	config, err := LoadCephCLIConfig()
	if err != nil {
		log.Panicf("failed to load config: %v", err)
	}

	rand.Shuffle(len(config.Hosts), func(i, j int) {
		config.Hosts[i], config.Hosts[j] = config.Hosts[j], config.Hosts[i]
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

	err = cephsetup.Setup(tempDir, config.Hosts, config.Keyring)
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
