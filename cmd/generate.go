package main

import (
	"context"

	"github.com/jolfzverb/codegen/internal/generator"
	"github.com/jolfzverb/codegen/internal/generator/options"
)

func main() {
	opts, err := options.GetOptions()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	gen := generator.NewGenerator(opts)
	err = gen.Generate(ctx)
	if err != nil {
		panic(err)
	}
}
