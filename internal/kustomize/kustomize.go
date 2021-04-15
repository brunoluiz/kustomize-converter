package kustomize

import (
	"regexp"
	"strings"

	"github.com/brunoluiz/kustomize-converter/internal/kustomize/inspect"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

type Kustomize struct {
	APIVersion       string            `yaml:"apiVersion" json:"apiVersion"`
	Kind             string            `yaml:"kind" json:"kind"`
	Namespace        string            `yaml:"namespace" json:"namespace"`
	Secrets          []ConfigGenerator `yaml:"secretGenerator,omitempty" json:"secretGenerator,namespace"`
	Configs          []ConfigGenerator `yaml:"configMapGenerator,omitempty" json:"configMapGenerator,omitempty"`
	Resources        []string          `yaml:"resources,omitempty" json:"resources,omitempty"`
	ResourcesData    map[string]string `yaml:"-" json:"-"`
	Processed        []string          `yaml:"-" json:"-"`
	deserializer     runtime.Decoder   `yaml:"-" json:"-"`
	generators       bool              `yaml:"-" json:"-"`
	processedLog     bool              `yaml:"-" json:"-"`
	baseFolder       string            `yaml:"-" json:"-"`
	secretsFolder    string            `yaml:"-" json:"-"`
	configsFolder    string            `yaml:"-" json:"-"`
	namespaceMatcher *regexp.Regexp    `yaml:"-" json:"-`
}

func (k *Kustomize) addProcessed(p string) {
	if !k.processedLog {
		return
	}
	k.Processed = append(k.Processed, p)
}

func (k *Kustomize) AddYAML(ypath string, data []byte) error {
	yamls, handleable, err := inspect.YAML(data)
	if err != nil {
		return err
	}

	// If it is a mixed manifest definitions (different .Kind tags) or generators
	// are disabled, don't do anything besides listing it in the resources list
	if !handleable || !k.generators {
		k.addResource(ypath, data)
		k.addProcessed(ypath)
		return nil
	}

	for _, y := range yamls {
		k.handle(ypath, y)
	}

	return nil
}

func (k *Kustomize) handle(path string, y []byte) {
	res, _, err := k.deserializer.Decode(y, nil, nil)
	// If there is a decode error, do not throw an error: just add it as a resource
	if err != nil {
		k.addResource(path, y)
		k.addProcessed(path)
		return
	}

	switch obj := res.(type) {
	case *corev1.Secret:
		k.Secrets = append(k.Secrets, newSecretGenerator(k.secretsFolder, path, obj))
	case *corev1.ConfigMap:
		k.Configs = append(k.Configs, newConfigMapGenerator(k.configsFolder, path, obj))
	default:
		k.addResource(path, y)
	}

	k.addProcessed(path)
}

func (k *Kustomize) addResource(ypath string, data []byte) {
	p := strings.ReplaceAll(ypath, k.baseFolder, ".")

	if _, ok := k.ResourcesData[p]; !ok {
		k.Resources = append(k.Resources, p)
		k.addProcessed(ypath)
	}

	k.ResourcesData[p] = k.namespaceMatcher.ReplaceAllString(string(data), "")
}

type KustomizeOption func(k *Kustomize)

func WithGenerators(enabled bool) func(k *Kustomize) {
	return func(k *Kustomize) {
		k.generators = enabled
	}
}

func WithBaseFolder(folder string) func(k *Kustomize) {
	return func(k *Kustomize) {
		k.baseFolder = folder
	}
}

func WithNamespace(ns string) func(k *Kustomize) {
	return func(k *Kustomize) {
		k.Namespace = ns
		k.namespaceMatcher = regexp.MustCompile(`.*(namespace:).*(` + ns + `).*\n`)
	}
}

func WithProcessedLog(c bool) func(k *Kustomize) {
	return func(k *Kustomize) {
		k.processedLog = c
	}
}

func WithConfigsFolder(f string) func(k *Kustomize) {
	return func(k *Kustomize) {
		k.configsFolder = f
	}
}

func WithSecretsFolder(f string) func(k *Kustomize) {
	return func(k *Kustomize) {
		k.secretsFolder = f
	}
}

func New(opts ...KustomizeOption) *Kustomize {
	k := &Kustomize{
		Kind:          "Kustomization",
		APIVersion:    "kustomize.config.k8s.io/v1beta1",
		Secrets:       []ConfigGenerator{},
		Configs:       []ConfigGenerator{},
		ResourcesData: map[string]string{},
		deserializer:  scheme.Codecs.UniversalDeserializer(),
		generators:    true,
		processedLog:  false,
		secretsFolder: "secrets",
		configsFolder: "configs",
	}

	for _, opt := range opts {
		opt(k)
	}
	return k
}
