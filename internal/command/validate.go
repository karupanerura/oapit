package command

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

type ValidateCommand struct {
	Schema  string `short:"s" name:"schema" help:"Schema component name" required:""`
	Request string ``
	As      string `name:"as" enum:"request,response" help:"validate as request/response" default:"request"`
	Payload string `arg:"" help:"Filepath for JSON payload" type:"existingfile" required:""`
}

func (c *ValidateCommand) Run(ctx context.Context, opts Options) error {
	var payload any
	if c.Payload == "-" {
		if err := json.NewDecoder(os.Stdin).Decode(&payload); err != nil {
			return fmt.Errorf("failed to decode JSON %q: %w", c.Payload, err)
		}
	} else {
		f, err := os.Open(c.Payload)
		if err != nil {
			return fmt.Errorf("failed to open file %q: %w", c.Payload, err)
		}
		defer f.Close()
		if err = json.NewDecoder(f).Decode(&payload); err != nil {
			return fmt.Errorf("failed to decode JSON %q: %w", c.Payload, err)
		}
	}

	root, err := opts.LoadSchema(ctx)
	if err != nil {
		return fmt.Errorf("failed to load schema %q: %w", opts.SchemaFile, err)
	}

	validationOpts := []openapi3.SchemaValidationOption{openapi3.EnableFormatValidation(), openapi3.MultiErrors()}
	switch c.As {
	case "request":
		validationOpts = append(validationOpts, openapi3.VisitAsRequest())
	case "response":
		validationOpts = append(validationOpts, openapi3.VisitAsResponse())
	}

	schema := root.Components.Schemas[c.Schema]
	if err := schema.Value.VisitJSON(payload, validationOpts...); err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	log.Println("No validation errors")
	return nil
}
