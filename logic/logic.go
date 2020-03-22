package logic

import (
	"fmt"
	"github.com/bene/graphcms/logic/converter"
	"github.com/bene/graphcms/logic/types"
	"github.com/graphql-go/graphql"
)

type Logic struct {
	fields graphql.Fields
	schema graphql.Schema
}

func (l *Logic) UpdateSchema() error {
	rootQuery := graphql.ObjectConfig{Name: "Query", Fields: l.fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}

	schema, err := graphql.NewSchema(schemaConfig)
	if err == nil {
		l.schema = schema
	}

	return err
}

func (l *Logic) CreateMask(mask types.Mask) error {

	field, err := converter.MaskToSingleField(mask)
	if err != nil {
		return err
	}
	l.fields[mask.Name] = &field

	// Add query all
	l.fields[fmt.Sprintf("%ss", mask.Name)] = &graphql.Field{
		Type: graphql.NewList(field.Type),
		Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
			return struct{}{}, nil
		},
	}

	// Add query single
	singleField := field
	singleField.Args = map[string]*graphql.ArgumentConfig{
		"id": {
			Type: graphql.NewNonNull(graphql.ID),
		},
	}

	// TODO: add connection query
	// l.fields[fmt.Sprintf("%ssConnection", mask.Name)]

	return l.UpdateSchema()
}

func (l *Logic) GetSchema() *graphql.Schema {
	return &l.schema
}

func NewLogic() (*Logic, error) {

	systemField := graphql.NewObject(graphql.ObjectConfig{
		Name: "System",
		Fields: graphql.Fields{
			"status": &graphql.Field{
				Description: "",
				Type:        graphql.String,
				Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
					return "unknown", nil
				},
			},
			"version": &graphql.Field{
				Description: "",
				Type:        graphql.String,
				Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
					return "0.0.0", nil
				},
			},
		},
	})

	fields := graphql.Fields{
		"_system": &graphql.Field{
			Type:        systemField,
			Description: "",
			Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
				return struct{}{}, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "Query", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, err
	}

	return &Logic{
		fields: fields,
		schema: schema,
	}, nil
}
