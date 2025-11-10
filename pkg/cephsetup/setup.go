package cephsetup

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var templates embed.FS

func Setup(outputDir string, hosts []string, key string) error {
	tmpl, err := template.ParseFS(templates, "templates/*.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse templates: %w", err)
	}

	cephConf := filepath.Join(outputDir, "ceph.conf")

	err = writeTemplate(tmpl, cephConf, "ceph.conf.tmpl", strings.Join(hosts, ","))
	if err != nil {
		return fmt.Errorf("failed to write ceph.conf: %w", err)
	}

	keyringFile := filepath.Join(outputDir, "ceph.client.admin.keyring")

	err = writeTemplate(tmpl, keyringFile, "ceph.client.admin.keyring.tmpl", key)
	if err != nil {
		return fmt.Errorf("failed to write keyring: %w", err)
	}

	return nil
}

func writeTemplate(tmpl *template.Template, path string, templateName string, data string) error {
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
