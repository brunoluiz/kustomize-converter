package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/brunoluiz/kustomize-converter/internal/kustomizer"
	"github.com/brunoluiz/kustomize-converter/internal/writer"
	"github.com/peterbourgon/ff/v3"
)

type Config struct {
	Folder           *string
	OutputFolder     *string
	EnableGenerators *bool
}

func main() {
	c := Config{}
	fs := flag.NewFlagSet("kustomize-converter", flag.ExitOnError)
	c.Folder = fs.String("folder", "", "kubernetes manifests folder")
	c.OutputFolder = fs.String("output-folder", "", "output folder (can be the same as --folder)")
	c.EnableGenerators = fs.Bool("generators", true, "toggle secret and configMapGenerator transforms")

	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarNoPrefix()); err != nil {
		log.Fatal(err)
	}

	if err := run(context.Background(), c); err != nil {
		fmt.Println(err)
	}
}

func run(ctx context.Context, c Config) error {
	k := kustomizer.FromFS(*c.Folder, *c.EnableGenerators)
	w := writer.FS{Folder: *c.OutputFolder}

	if err := k.ParseYAML(); err != nil {
		return err
	}

	return w.Write(k.Output)
}
