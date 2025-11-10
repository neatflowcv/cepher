package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"text/template"

	"github.com/neatflowcv/cepher/pkg/cephcli"
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

//go:embed templates/*.tmpl
var templates embed.FS

func main() {
	log.Println("version", version())

	tmpl, err := template.ParseFS(templates, "templates/*.tmpl")
	if err != nil {
		log.Panicf("failed to parse templates: %v", err)
	}

	hosts := os.Getenv("CEPH_MON_HOSTS")
	if hosts == "" {
		log.Panicf("CEPH_MON_HOSTS must be set")
	}

	monHosts := strings.Split(hosts, ",")
	rand.Shuffle(len(monHosts), func(i, j int) {
		monHosts[i], monHosts[j] = monHosts[j], monHosts[i]
	})

	keyring := os.Getenv("CEPH_KEYRING")
	if keyring == "" {
		log.Panicf("CEPH_KEYRING must be set")
	}

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

	cephConf := filepath.Join(tempDir, "ceph.conf")

	err = writeTemplate(tmpl, cephConf, "ceph.conf.tmpl", strings.Join(monHosts, ","))
	if err != nil {
		log.Panicf("failed to write ceph.conf: %v", err)
	}

	keyringFile := filepath.Join(tempDir, "ceph.client.admin.keyring")

	err = writeTemplate(tmpl, keyringFile, "ceph.client.admin.keyring.tmpl", keyring)
	if err != nil {
		log.Panicf("failed to write keyring: %v", err)
	}

	client := cephcli.NewClient(tempDir, "20.1.1")

	health, err := client.GetHealth(context.Background())
	if err != nil {
		log.Panicf("failed to get health detail from ceph: %v", err)
	}

	log.Println("health", health)
}

func writeTemplate(tmpl *template.Template, path string, templateName string, data any) error {
	file, err := os.Create(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	err = tmpl.ExecuteTemplate(file, templateName, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
