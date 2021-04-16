package kustomize

import (
	"fmt"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

// ConfigGenerator Defines an instance of Kustomize ConfigGenerators (secret and configMap generator).
// Contains not only the YAML/JSON fields, but as well the file and env data and helper methods.
type ConfigGenerator struct {
	Name  string   `yaml:"name" json:"name"`
	Type  string   `yaml:"type,omitempty" json:"type,omitempty"`
	Envs  []string `yaml:"envs,omitempty" json:"envs,omitempty"`
	Files []string `yaml:"files,omitempty" json:"files,omitempty"`

	Path     string            `yaml:"-" json:"-"`
	FileDir  string            `yaml:"-" json:"-"`
	EnvData  map[string]string `yaml:"-" json:"-"`
	FileData map[string]string `yaml:"-" json:"-"`
}

// Add Adds a config to the ConfigGenerator. If it is a multi-line content, it will be added to .Files
// and .FileData. Otherwise, it will be added to .Env and .EnvData.
func (c *ConfigGenerator) Add(k, content string) {
	// If a file is multi-line, probably it is a file. Hence, add it to
	// .Files and .FilesData, instead of .Env and .EnvsData
	if strings.Contains(content, "\n") {
		file := c.Name + "-" + strings.TrimPrefix(k, ".")
		p := filepath.Join(c.Path, file)
		c.FileData[p] = content
		c.Files = append(c.Files, k+"="+p)

		return
	}

	c.EnvData[k] = content
	if len(c.Envs) == 0 {
		c.Envs = []string{filepath.Join(c.Path, c.Name)}
	}
}

// GetEnv Parses and returns all env data in a nicely formatted string.
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
