package main

import (
	"flag"
	"log"
	"os"

	"github.com/brunoluiz/kustomize-converter/internal/kustomize"
	"github.com/brunoluiz/kustomize-converter/internal/loader"
	"github.com/brunoluiz/kustomize-converter/internal/writer"
	"github.com/peterbourgon/ff/v3"
)

func main() {
	fs := flag.NewFlagSet("kustomize-converter", flag.ExitOnError)

	var (
		folder           = fs.String("folder", "", "kubernetes manifests folder")
		outputFolder     = fs.String("output-folder", "", "output folder (can be the same as --folder)")
		enableGenerators = fs.Bool("generators", true, "toggle secret and configMapGenerator transforms")
		namespace        = fs.String("namespace", "", "set a kubernetes namespace, instead of trying to infer from files")
		clean            = fs.Bool("clean", false, "if set, it will clear up resources from 'folder' before generating new")
	)

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix()); err != nil {
		log.Fatal(err)
	}

	k, err := loader.FromFS(*folder,
		kustomize.WithBaseFolder(*folder),
		kustomize.WithGenerators(*enableGenerators),
		kustomize.WithNamespace(*namespace),
		kustomize.WithProcessedLog(*clean),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := writer.ToFS(*outputFolder, *clean).Write(k); err != nil {
		log.Fatal(err)
	}
}
