package command

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
)

type ValidateCommand struct {
	Schema  ValidateSchemaCommand  `cmd:""`
	Payload ValidatePayloadCommand `cmd:""`
}

type ValidateSchemaCommand struct {
	AllowExtraSiblingFields []string `name:"allow-extra-sibling-fields" optional:"" help:"TBD"`
	ExamplesValidation      bool     `name:"validate-examples" negatable:"" default:"true" help:"TBD"`
	DefaultsValidation      bool     `name:"validate-defaults" negatable:"" default:"true" help:"TBD"`
	PatternValidation       bool     `name:"validate-patterns" negatable:"" default:"true" help:"TBD"`
}

func (c *ValidateSchemaCommand) validationOptions() (opts []openapi3.ValidationOption) {
	if len(c.AllowExtraSiblingFields) != 0 {
		opts = append(opts, openapi3.AllowExtraSiblingFields(c.AllowExtraSiblingFields...))
	}

	// validate examples
	if c.ExamplesValidation {
		opts = append(opts, openapi3.EnableExamplesValidation())
	} else {
		opts = append(opts, openapi3.DisableExamplesValidation())
	}

	// validate defaults
	if c.DefaultsValidation {
		opts = append(opts, openapi3.EnableSchemaDefaultsValidation())
	} else {
		opts = append(opts, openapi3.DisableSchemaDefaultsValidation())
	}

	// validate pattern
	if c.PatternValidation {
		opts = append(opts, openapi3.EnableSchemaPatternValidation())
	} else {
		opts = append(opts, openapi3.DisableSchemaPatternValidation())
	}

	return
}

func (c *ValidateSchemaCommand) Run(ctx context.Context, opts Options) error {
	_, err := opts.LoadSchema(ctx, c.validationOptions()...)
	if err != nil {
		return fmt.Errorf("failed to load schema %q: %w", opts.SchemaFile, err)
	}

	log.Println("No validation errors")
	return nil
}

type ValidatePayloadCommand struct {
	Schema  string `short:"s" name:"schema" help:"Schema component name" required:""`
	As      string `name:"as" enum:"request,response" help:"validate as request/response" default:"request"`
	Payload string `arg:"" help:"Filepath for JSON payload" type:"existingfile" required:""`
}

func (c *ValidatePayloadCommand) Run(ctx context.Context, opts Options) error {
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
		var mErr openapi3.MultiError
		if errors.As(err, &mErr) {
			for _, err := range mErr {
				log.Println(err.Error())
			}
			return fmt.Errorf("validation errors: %d errors", len(mErr))
		}

		var sErr *openapi3.SchemaError
		if errors.As(err, &sErr) {
			log.Println(sErr.Error())
			return fmt.Errorf("validation error: 1 error")
		}

		return fmt.Errorf("unknown validation error: %w", err)
	}

	log.Println("No validation errors")
	return nil
}
