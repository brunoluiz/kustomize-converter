package inspect

import (
	"bytes"
	"strings"

	"gopkg.in/yaml.v2"
)

func YAML(data []byte) (out [][]byte, mixed bool, err error) {
	kinds := map[string]bool{}
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

		kinds[kind] = true

		out = append(out, y)
	}

	return out, len(kinds) > 1, nil
}
