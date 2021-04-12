## kustomize-converter

Converts Kubernetes YAML files to Kustomize

## Running

```
Usage of kustomize-converter:
  -clean
        if set to true, it will clear up resources from source folder before generating output
  -configs-folder string
        which folder should the ConfigMaps be placed in the output folder (default "configs")
  -folder string
        kubernetes manifests source folder
  -generators
        if set to false, disable secret and configMapGenerator transforms (default true)
  -namespace string
        set a kubernetes namespace, instead of trying to infer from files
  -output-folder string
        kubernetes manifest output folder (can be the same as --folder)
  -secrets-folder string
        which folder should the Secrets be placed in the output folder (default "secrets")
```

## Output details

- Files with mixed types of resources will be added to `$.resources` and will not be transformed in any way
- Files added to `$.resources` will have the property `$.namespace` removed, as it is defined in the `kustomization.yaml`
- Files with one or more `Secret` resources will be added to `$.secretGenerator` and will be transformed
- Files with one or more `ConfigMap` resources will be added to `$.configMapGenerator` and will be transformed
- `ConfigMap` and `Secret` with multi-line entries will be exported as a file to `$.[generator].files`
- `ConfigMap` and `Secret` with single-line entries will be exported as an env to `$.[generator].envs`
- Transformed `ConfigMap` resources will be placed at `./configs/${ prefix }-${ obj.name }`
- Transformed `Secret` resources will be placed at `./secrets/${ prefix }-${ obj.name }`

### Example output file and folder structure

```
# resources which are neither secrets or config maps
- api/
  - deployment.yaml
  - reader.yaml
- namespace.yml 
- auth.yml 
- database.yml 

# files used by the config and secrets generator
- secrets/
  - api
  - database
  - clients.json
- configs/
  - api

# contain all kustomize configs
- kustomization.yml
```
