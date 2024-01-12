package shared

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

var (
	ErrSchemaNotFound = errors.New("schema not found")
	ErrInvalidSchema  = errors.New("invalid schema")
)

type SchemaValidator struct {
	fs      embed.FS
	schemas map[string]*jsonschema.Schema
}

func NewSchemaValidator(fs embed.FS) *SchemaValidator {
	return &SchemaValidator{fs, make(map[string]*jsonschema.Schema)}
}

func (v *SchemaValidator) ValidateBytes(schemaName string, docBytes []byte) error {
	var doc any
	if err := json.Unmarshal(docBytes, &doc); err != nil {
		return err
	}

	return v.Validate(schemaName, doc)
}

func (v *SchemaValidator) Validate(schemaName string, doc any) error {
	schema, err := v.getSchema(schemaName)
	if err != nil {
		return err
	}

	if err = schema.Validate(doc); err != nil {
		var errs []string
		if verr, ok := err.(*jsonschema.ValidationError); ok {
			for _, e := range verr.BasicOutput().Errors {
				if e.KeywordLocation == "" || e.Error == "oneOf failed" || e.Error == "allOf failed" {
					continue
				}

				if e.InstanceLocation == "" {
					errs = append(errs, e.Error)
				} else {
					errs = append(errs, fmt.Sprintf("%s: %s", e.InstanceLocation, e.Error))
				}
			}
			return fmt.Errorf("%v", errs)
		} else {
			return fmt.Errorf("could not validate: %s", err)
		}
	}

	return nil
}

func (v *SchemaValidator) getSchema(name string) (*jsonschema.Schema, error) {
	schema, found := v.schemas[name]

	if !found {
		s, err := v.fs.ReadFile(name)
		if err != nil {
			return nil, ErrSchemaNotFound
		}

		compiler := jsonschema.NewCompiler()
		compiler.Draft = jsonschema.Draft7

		if err := compiler.AddResource("schema.json", strings.NewReader(string(s))); err != nil {
			return nil, ErrInvalidSchema
		}

		schema, err = compiler.Compile("schema.json")
		if err != nil {
			return nil, ErrInvalidSchema
		}

		v.schemas[name] = schema
	}

	return schema, nil
}
