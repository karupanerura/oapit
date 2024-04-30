package command

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
)

// Options is a global options
type Options struct {
	SchemaFile       string `short:"f" name:"schema-file" type:"existingfile" help:"OpenAPI 3.0 schema" required:""`
	ResolveReference bool   `name:"resolve-reference" negatable:"" default:"true" help:"TBD"`
}

func (o *Options) LoadSchema(ctx context.Context, opts ...openapi3.ValidationOption) (*openapi3.T, error) {
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

	if len(opts) != 0 {
		if err = doc.Validate(loader.Context, opts...); err != nil {
			return nil, err
		}
	}

	return doc, nil
}
