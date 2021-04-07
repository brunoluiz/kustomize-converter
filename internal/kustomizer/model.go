package kustomizer

type Kustomize struct {
	APIVersion    string            `yaml:"apiVersion"`
	Kind          string            `yaml:"kind"`
	Namespace     string            `yaml:"namespace"`
	Secrets       []ConfigGenerator `yaml:"secretGenerator,omitempty"`
	Configs       []ConfigGenerator `yaml:"configMapGenerator,omitempty"`
	Resources     []string          `yaml:"resources,omitempty"`
	ResourcesData map[string]string `yaml:"-"`
}

func (k *Kustomize) SetNamespace(n string) {
	if k.Namespace != "" {
		return
	}
	k.Namespace = n
}

func NewKustomize() *Kustomize {
	return &Kustomize{
		Kind:          "Kustomization",
		APIVersion:    "kustomize.config.k8s.io/v1beta1",
		Secrets:       []ConfigGenerator{},
		Configs:       []ConfigGenerator{},
		ResourcesData: map[string]string{},
	}
}
