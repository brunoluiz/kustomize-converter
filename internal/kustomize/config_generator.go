package kustomize

import (
	"fmt"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

type ConfigGenerator struct {
	Name     string            `yaml:"name" json:"name"`
	Type     string            `yaml:"type,omitempty" json:"type,omitempty"`
	Envs     []string          `yaml:"envs,omitempty" json:"envs,omitempty"`
	Files    []string          `yaml:"files,omitempty" json:"files,omitempty"`
	Path     string            `yaml:"-" json:"-"`
	FileDir  string            `yaml:"-" json:"-"`
	EnvData  map[string]string `yaml:"-" json:"-"`
	FileData map[string]string `yaml:"-" json:"-"`
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

func newSecretGenerator(baseFolder, ypath string, obj *corev1.Secret) ConfigGenerator {
	secret := ConfigGenerator{
		Name:     obj.Name,
		FileDir:  filepath.Dir(ypath),
		Path:     baseFolder,
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

	// It is always Opaque by default, but otherwise we might want to change it
	if obj.Type != corev1.SecretTypeOpaque {
		secret.Type = string(obj.Type)
	}

	return secret
}

func newConfigMapGenerator(baseFolder, ypath string, obj *corev1.ConfigMap) ConfigGenerator {
	config := ConfigGenerator{
		Name:     obj.Name,
		FileDir:  filepath.Dir(ypath),
		Path:     baseFolder,
		Envs:     []string{},
		EnvData:  map[string]string{},
		FileData: map[string]string{},
	}

	for k, v := range obj.Data {
		config.Add(k, v)
	}

	return config
}
