package writer

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/brunoluiz/kustomize-converter/internal/kustomize"
	"gopkg.in/yaml.v2"
)

type FS struct {
	folder         string
	clearProcessed bool
}

func ToFS(folder string, clearProcessed bool) FS {
	return FS{folder: folder, clearProcessed: clearProcessed}
}

func (w FS) Write(k *kustomize.Kustomize) error {
	if err := w.writeGenerators(k.Secrets); err != nil {
		return err
	}

	if err := w.writeGenerators(k.Configs); err != nil {
		return err
	}

	if err := w.writeResources(k.ResourcesData); err != nil {
		return err
	}

	if err := w.deleteResources(k.Processed); err != nil {
		return err
	}

	d, err := yaml.Marshal(&k)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(w.folder, "kustomization.yaml"), d, 0600)
}

func (w FS) writeGenerators(configs []kustomize.ConfigGenerator) error {
	for _, config := range configs {
		p := filepath.Join(w.folder, config.Path)
		if err := os.MkdirAll(p, 0755); err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}

		if len(config.Envs) > 0 {
			p := filepath.Join(w.folder, config.Envs[0])
			if err := ioutil.WriteFile(p, []byte(config.GetEnv()), 0600); err != nil {
				return err
			}
		}

		for fpath, file := range config.FileData {
			p := filepath.Join(w.folder, fpath)
			if err := ioutil.WriteFile(p, []byte(file), 0600); err != nil {
				return err
			}
		}
	}

	return nil
}

func (w FS) writeResources(resources map[string]string) error {
	for fres, res := range resources {
		p := filepath.Join(w.folder, filepath.Dir(fres))
		if err := os.MkdirAll(p, 0755); err != nil && !errors.Is(err, os.ErrExist) {
			return err
		}

		p = filepath.Join(w.folder, fres)
		if err := ioutil.WriteFile(p, []byte(res), 0600); err != nil {
			return err
		}
	}

	return nil
}

func (w FS) deleteResources(resources []string) error {
	if !w.clearProcessed {
		return nil
	}

	for _, res := range resources {
		if err := os.Remove(res); err != nil {
			return err
		}
	}

	return nil
}
