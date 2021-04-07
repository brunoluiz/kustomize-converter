package kustomizer

import (
	"fmt"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

type ConfigGenerator struct {
	Name     string            `yaml:"name"`
	Envs     []string          `yaml:"envs,omitempty"`
	Files    []string          `yaml:"files,omitempty"`
	Path     string            `yaml:"-"`
	FileDir  string            `yaml:"-"`
	EnvData  map[string]string `yaml:"-"`
	FileData map[string]string `yaml:"-"`
}

func (c *ConfigGenerator) Add(k, v string) {
	if strings.Contains(v, "\n") {
		p := filepath.Join(c.Path, c.Name+"-"+k)
		c.FileData[p] = v
		c.Files = append(c.Files, k+"="+p)

		return
	}

	c.EnvData[k] = v
	if len(c.Envs) == 0 {
		c.Envs = []string{filepath.Join(c.Path, c.Name)}
	}
}

func (c *ConfigGenerator) GetEnv() (s string) {
	for k, v := range c.EnvData {
		s += fmt.Sprintf("%s=%s\n", k, strings.ReplaceAll(v, "\n", `\n`))
	}

	return s
}

func newSecretGenerator(ypath string, obj *corev1.Secret) ConfigGenerator {
	secret := ConfigGenerator{
		Name:     obj.Name,
		FileDir:  filepath.Dir(ypath),
		Path:     "secrets",
		Envs:     []string{},
		EnvData:  map[string]string{},
		FileData: map[string]string{},
	}

	for k, v := range obj.Data {
		secret.Add(k, string(v))
	}

	for k, v := range obj.StringData {
		secret.Add(k, v)
	}

	return secret
}

func newConfigMapGenerator(ypath string, obj *corev1.ConfigMap) ConfigGenerator {
	config := ConfigGenerator{
		Name:     obj.Name,
		FileDir:  filepath.Dir(ypath),
		Path:     "configs",
		Envs:     []string{},
		EnvData:  map[string]string{},
		FileData: map[string]string{},
	}

	for k, v := range obj.Data {
		config.Add(k, v)
	}

	return config
}
