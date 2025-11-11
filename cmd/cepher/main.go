package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/neatflowcv/cepher/internal/app/flow"
	"github.com/neatflowcv/cepher/internal/pkg/client/core"
	"github.com/neatflowcv/cepher/internal/pkg/idgenerator/ulid"
	"github.com/neatflowcv/cepher/internal/pkg/repository/file"
	"github.com/neatflowcv/cepher/pkg/cephrest"
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

	const (
		timeout = 10 * time.Second
	)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Panicf("failed to get home directory: %v", err)
	}

	storagePath := filepath.Join(homeDir, ".local/share/cepher")

	const permission = 0750

	err = os.MkdirAll(storagePath, permission)
	if err != nil {
		log.Fatalf("failed to create storage directory: %v", err)
	}

	repository := file.NewRepository(storagePath)
	service := flow.NewService(ulid.NewGenerator(), core.NewFactory(), repository)

	handler, err := NewHandler(service)
	if err != nil {
		log.Panicf("failed to create handler: %v", err)
	}
	defer handler.Close()

	server := &http.Server{ //nolint:exhaustruct
		ReadHeaderTimeout: timeout,
		Addr:              ":8080",
		Handler:           handler.Get(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panicf("failed to start server: %v", err)
	}
}
