package generator

import (
	"context"
	"io"
	"log/slog"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-faster/errors"
	"github.com/jolfzverb/codegen/internal/generator/options"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const directoryPermissions = 0o755

type Generator struct {
	Opts *options.Options

	SchemasFile  *SchemasFile
	HandlersFile *HandlersFile
	yaml         *openapi3.T

	// strings
	PackageName      string
	ImportPrefix     string
	ModelsImportPath string
	CurrentYAMLFile  string

	YAMLFilesToProcess []string
	YAMLFilesProcessed map[string]bool
}

func NewGenerator(opts *options.Options) *Generator {
	return &Generator{
		Opts:               opts,
		YAMLFilesToProcess: opts.YAMLFiles,
		YAMLFilesProcessed: make(map[string]bool),
	}
}

func (g *Generator) WriteToOutput(modelsOutput io.Writer, handlersOutput io.Writer) error {
	const op = "generator.Generator.WriteToOutput"
	err := g.WriteSchemasToOutput(modelsOutput)
	if err != nil {
		return errors.Wrap(err, op)
	}
	err = g.WriteHandlersToOutput(handlersOutput)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (g *Generator) Gen() {
	const op = "generator.Generate"

	// one time
	g.InitHandlerFields(g.PackageName)

	if g.yaml.Paths != nil && len(g.yaml.Paths.Map()) > 0 {
		err := g.ProcessPaths(g.yaml.Paths)
		if err != nil {
			panic(errors.Wrap(err, op))
		}
	}

	if g.yaml.Components != nil && g.yaml.Components.Schemas != nil {
		err := g.ProcessSchemas(g.yaml.Components.Schemas)
		if err != nil {
			panic(errors.Wrap(err, op))
		}
	}

}

func (g *Generator) GetModelName(yamlFilePath string) string {
	parts := strings.Split(yamlFilePath, "/")
	if len(parts) == 0 {
		return ""
	}
	fileName := parts[len(parts)-1]
	fileName = strings.TrimSuffix(fileName, ".yaml")
	fileName = strings.TrimSuffix(fileName, ".yml")
	fileName = strings.ReplaceAll(fileName, "_", "")
	fileName = strings.ReplaceAll(fileName, "-", "")

	lowerCaser := cases.Lower(language.Und)

	return lowerCaser.String(fileName)
}

func (g *Generator) PrepareAndRead(reader io.Reader) error {
	const op = "generator.PrepareAndRead"
	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	var err error
	data, err := io.ReadAll(reader)
	if err != nil {
		return errors.Wrap(err, op)
	}
	url, err := url.Parse(g.CurrentYAMLFile)
	if err != nil {
		return errors.Wrap(err, op)
	}
	g.yaml, err = loader.LoadFromDataWithPath(data, url)
	if err != nil {
		return errors.Wrap(err, op)
	}
	err = g.yaml.Validate(ctx)
	if err != nil {
		return errors.Wrap(err, op)
	}

	g.NewSchemasFile()
	g.NewHandlersFile()

	return nil
}

func (g *Generator) PrepareFiles() error {
	const op = "generator.PrepareFiles"

	file, err := os.Open(g.CurrentYAMLFile)
	if err != nil {
		return errors.Wrap(err, op)
	}
	defer file.Close()

	reader := io.Reader(file)

	g.PackageName = g.GetModelName(g.CurrentYAMLFile)

	handlersPath := path.Join(g.Opts.DirPrefix, "generated", g.PackageName)
	schemasPath := path.Join(handlersPath, g.GetCurrentModelsPackage())
	err = os.MkdirAll(schemasPath, directoryPermissions)
	if err != nil {
		return errors.Wrap(err, op)
	}

	g.ImportPrefix = path.Join(g.Opts.PackagePrefix, "generated", g.PackageName)
	g.ModelsImportPath = path.Join(g.ImportPrefix, g.GetCurrentModelsPackage())
	err = g.PrepareAndRead(reader)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (g *Generator) GenerateFiles() error {
	g.Gen()
	return nil
}
func (g *Generator) WriteOutFiles() error {
	const op = "generator.WriteOutFiles"

	handlersPath := path.Join(g.Opts.DirPrefix, "generated", g.PackageName)
	schemasPath := path.Join(handlersPath, g.GetCurrentModelsPackage())
	schemasOutput, err := os.Create(path.Join(schemasPath, "models.go"))
	if err != nil {
		return errors.Wrap(err, op)
	}
	defer schemasOutput.Close()

	handlersOutput, err := os.Create(path.Join(handlersPath, "handlers.go"))
	if err != nil {
		return errors.Wrap(err, op)
	}
	defer handlersOutput.Close()

	err = g.WriteToOutput(schemasOutput, handlersOutput)
	if err != nil {
		return errors.Wrap(err, op)
	}
	return nil
}

func (g *Generator) Generate(ctx context.Context) error {
	const op = "generator.Generate"

	for len(g.YAMLFilesToProcess) > 0 {
		g.CurrentYAMLFile = g.YAMLFilesToProcess[0]

		if g.YAMLFilesProcessed[g.CurrentYAMLFile] {
			g.YAMLFilesToProcess = g.YAMLFilesToProcess[1:]
			continue
		}
		slog.Info("Processing file", "file", g.CurrentYAMLFile)

		err := g.PrepareFiles()
		if err != nil {
			return errors.Wrap(err, op)
		}
		err = g.GenerateFiles()
		if err != nil {
			return errors.Wrap(err, op)
		}
		err = g.WriteOutFiles()
		if err != nil {
			return errors.Wrap(err, op)
		}
		g.YAMLFilesProcessed[g.CurrentYAMLFile] = true
		g.YAMLFilesToProcess = g.YAMLFilesToProcess[1:]
	}

	return nil
}
