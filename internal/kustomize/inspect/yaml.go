package inspect

import (
	"bytes"
	"strings"

	"gopkg.in/yaml.v2"
)

var supportedTypes = map[string]bool{"Secret": true, "ConfigMap": true}

// YAML Return useful data about the YAML, which might decide how it will be
// parsed by the Kustomize instance
func YAML(data []byte) (out [][]byte, handleable bool, err error) {
	yamls := bytes.Split(data, []byte("\n---"))

	for _, y := range yamls {
		y := y

		// empty yaml
		var str string
		str = strings.ReplaceAll(string(y), "\n", "")
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

		if _, ok := supportedTypes[kind]; !ok {
			return yamls, false, nil
		}
	}

	return yamls, true, nil
}
