//nolint: lll
package main

import (
	"fmt"
	"os"

	"github.com/brunoluiz/kustomize-converter/internal/kustomize"
	"github.com/brunoluiz/kustomize-converter/internal/loader"
	"github.com/brunoluiz/kustomize-converter/internal/writer"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		&cli.StringFlag{Name: "folder", Usage: "kubernetes manifests source folder", Required: true},
		&cli.StringFlag{Name: "output-folder", Usage: "kubernetes manifest output folder (can be the same as --folder)", Required: true},
		&cli.StringFlag{Name: "namespace", Usage: "kubernetes namespace for this resources", Required: true},
		&cli.StringFlag{Name: "secrets-folder", Usage: "which folder should the ConfigMaps be placed in the output folder", Value: "secrets"},
		&cli.StringFlag{Name: "configs-folder", Usage: "which folder should the Secrets be placed in the output folder", Value: "configs"},
		&cli.BoolFlag{Name: "clean", Usage: "if set to true, it will clear up resources from source folder before generating output", Value: false},
		&cli.BoolFlag{Name: "generators", Usage: "if set to false, disable secret and configMapGenerator transforms", Value: true},
	}

	app.Action = func(c *cli.Context) error {
		k, err := loader.FromFS(c.String("folder"),
			kustomize.WithNamespace(c.String("namespace")),
			kustomize.WithConfigsFolder(c.String("configs-folder")),
			kustomize.WithSecretsFolder(c.String("secrets-folder")),
			kustomize.WithGenerators(c.Bool("generators")),
			kustomize.WithProcessedLog(c.Bool("clean")),
		)
		if err != nil {
			return errors.Wrap(err, "issue on loading kubernetes manifests")
		}

		if err := writer.ToFS(c.String("output-folder"), c.Bool("clean")).Write(k); err != nil {
			return errors.Wrap(err, "issue on writing kubernetes manifests")
		}

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err.Error())
	}
}
