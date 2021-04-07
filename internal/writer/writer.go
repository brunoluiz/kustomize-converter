package writer

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/brunoluiz/kustomize-converter/internal/kustomizer"
	"gopkg.in/yaml.v2"
)

type FS struct {
	Folder string
}

func (w *FS) Write(k *kustomizer.Kustomize) error {
	if err := w.writeGenerators(k.Secrets); err != nil {
		return err
	}

	if err := w.writeGenerators(k.Configs); err != nil {
		return err
	}

	if err := w.writeResources(k.ResourcesData); err != nil {
		return err
	}

	d, err := yaml.Marshal(&k)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(w.Folder, "kustomization.yml"), d, 0600)
}

func (w *FS) writeGenerators(configs []kustomizer.ConfigGenerator) error {
	for _, config := range configs {
		p := filepath.Join(w.Folder, config.Path)
		if err := os.MkdirAll(p, 0755); err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}

		if len(config.Envs) > 0 {
			p := filepath.Join(w.Folder, config.Envs[0])
			if err := ioutil.WriteFile(p, []byte(config.GetEnv()), 0644); err != nil {
				return err
			}
		}

		for fpath, file := range config.FileData {
			p := filepath.Join(w.Folder, fpath)
			if err := ioutil.WriteFile(p, []byte(file), 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *FS) writeResources(resources map[string]string) error {
	for fres, res := range resources {
		p := filepath.Join(w.Folder, filepath.Dir(fres))
		if err := os.MkdirAll(p, 0755); err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}

		p = filepath.Join(w.Folder, fres)
		if err := ioutil.WriteFile(p, []byte(res), 0644); err != nil {
			return err
		}
	}
	return nil
}
