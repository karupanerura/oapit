package command

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
)

// Options is a global options
type Options struct {
	SchemaFile              string   `short:"f" name:"schema-file" type:"existingfile" help:"OpenAPI 3.0 schema" required:""`
	ValidateSchema          bool     `name:"validate-schema" negatable:"" default:"true" help:"TBD"`
	ResolveReference        bool     `name:"resolve-reference" negatable:"" default:"true" help:"TBD"`
	AllowExtraSiblingFields []string `name:"allow-extra-sibling-fields" optional:"" help:"TBD"`
	ExamplesValidation      bool     `name:"validate-examples" negatable:"" default:"true" help:"TBD"`
	DefaultsValidation      bool     `name:"validate-defaults" negatable:"" default:"true" help:"TBD"`
	PatternValidation       bool     `name:"validate-patterns" negatable:"" default:"true" help:"TBD"`
}

func (o *Options) validationOptions() (opts []openapi3.ValidationOption) {
	if len(o.AllowExtraSiblingFields) != 0 {
		opts = append(opts, openapi3.AllowExtraSiblingFields(o.AllowExtraSiblingFields...))
	}

	// validate examples
	if o.ExamplesValidation {
		opts = append(opts, openapi3.EnableExamplesValidation())
	} else {
		opts = append(opts, openapi3.DisableExamplesValidation())
	}

	// validate defaults
	if o.DefaultsValidation {
		opts = append(opts, openapi3.EnableSchemaDefaultsValidation())
	} else {
		opts = append(opts, openapi3.DisableSchemaDefaultsValidation())
	}

	// validate pattern
	if o.PatternValidation {
		opts = append(opts, openapi3.EnableSchemaPatternValidation())
	} else {
		opts = append(opts, openapi3.DisableSchemaPatternValidation())
	}

	return
}

func (o *Options) LoadSchema(ctx context.Context) (*openapi3.T, error) {
	loader := openapi3.NewLoader()
	loader.Context = ctx
	if o.ResolveReference {
		loader.IsExternalRefsAllowed = true
	}

	var doc *openapi3.T
	var err error
	if o.SchemaFile == "-" {
		doc, err = loader.LoadFromStdin()
	} else {
		doc, err = loader.LoadFromFile(o.SchemaFile)
	}
	if err != nil {
		return nil, err
	}

	if o.ValidateSchema {
		if err = doc.Validate(loader.Context, o.validationOptions()...); err != nil {
			return nil, err
		}
	}

	return doc, nil
}
