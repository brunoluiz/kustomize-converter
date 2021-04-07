package kustomizer

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

type Kustomizer struct {
	Folder       string
	Output       *Kustomize
	Deserializer runtime.Decoder
}

func FromFS(folder string) *Kustomizer {
	return &Kustomizer{
		Folder:       folder,
		Output:       NewKustomize(),
		Deserializer: scheme.Codecs.UniversalDeserializer(),
	}
}

var namespaceMatcher = regexp.MustCompile(`(namespace:).*\n`)

func (kzr *Kustomizer) ParseYAML() error {
	return filepath.Walk(kzr.Folder, func(path string, info os.FileInfo, e error) error {
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

		return kzr.AddYAML(path, yaml)
	})
}

func (kzr *Kustomizer) AddYAML(ypath string, data []byte) error {
	yamls, mixed, err := kzr.inspect(data)
	if err != nil {
		return err
	}

	// If it is a mixed manifest definitions (different .Kind tags), don't do anything
	if mixed {
		kzr.addResource(ypath, data)
		return nil
	}

	for _, y := range yamls {
		res, _, err := kzr.Deserializer.Decode(y, nil, nil)
		// If there is a decode error, do not throw an error: just add it as a resource
		if err != nil {
			kzr.addResource(ypath, data)
			continue
		}

		switch obj := res.(type) {
		case *corev1.Secret:
			kzr.Output.Secrets = append(kzr.Output.Secrets, newSecretGenerator(ypath, obj))
			kzr.Output.SetNamespace(obj.Namespace)
		case *corev1.ConfigMap:
			kzr.Output.Configs = append(kzr.Output.Configs, newConfigMapGenerator(ypath, obj))
			kzr.Output.SetNamespace(obj.Namespace)
		default:
			kzr.addResource(ypath, data)
		}
	}

	return nil
}

func (kzr *Kustomizer) addResource(ypath string, data []byte) {
	p := strings.ReplaceAll(ypath, kzr.Folder, ".")
	kzr.Output.Resources = append(kzr.Output.Resources, p)
	d := string(data)
	d = namespaceMatcher.ReplaceAllString(d, "")
	kzr.Output.ResourcesData[p] = d
}

func (kzr *Kustomizer) inspect(data []byte) (out [][]byte, mixed bool, err error) {
	kinds := map[string]bool{}
	yamls := bytes.Split(data, []byte("\n---"))

	for _, y := range yamls {

		y := y
		// empty yaml
		str := strings.ReplaceAll(string(y), "\n", "")
		str = strings.Trim(str, " ")
		if str == "" {
			continue
		}

		info := map[string]interface{}{}
		if err := yaml.Unmarshal(y, &info); err != nil {
			return nil, false, err
		}

		kind, ok := info["kind"].(string)
		if !ok || kind == "" {
			continue
		}

		kinds[kind] = true
		out = append(out, y)
	}

	return out, len(kinds) > 1, nil
}
