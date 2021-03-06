package loader

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/brunoluiz/kustomize-converter/internal/kustomize"
)

// FromFS Load manifests from file system and generate a Kustomize instance
func FromFS(folder string, opts ...kustomize.KustomizeOption) (*kustomize.Kustomize, error) {
	k := kustomize.New(folder, opts...)

	err := filepath.Walk(folder, func(path string, info os.FileInfo, e error) error {
		if strings.Contains(path, "kustomization.") {
			return errors.New("manifests are already using kustomize")
		}

		// Ignore files which are not YAML
		if !(strings.Contains(path, "yml") || strings.Contains(path, "yaml")) {
			return nil
		}

		yaml, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		return k.AddYAML(path, yaml)
	})

	return k, err
}
