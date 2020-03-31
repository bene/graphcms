package logic

import (
	"fmt"
	"github.com/bene/graphcms/logic/converter"
	"github.com/bene/graphcms/logic/types"
	"github.com/graphql-go/graphql"
)

type Logic struct {
	types          []graphql.Type
	queryFields    graphql.Fields
	mutationFields graphql.Fields
	schema         graphql.Schema
}

func (l *Logic) UpdateSchema() error {
	query := graphql.ObjectConfig{Name: "Query", Fields: l.queryFields}
	mutation := graphql.ObjectConfig{Name: "Mutation", Fields: l.mutationFields}
	subscription := graphql.ObjectConfig{Name: "Subscription", Fields: l.mutationFields}
	schemaConfig := graphql.SchemaConfig{
		Query:        graphql.NewObject(query),
		Mutation:     graphql.NewObject(mutation),
		Subscription: graphql.NewObject(subscription),
	}

	schema, err := graphql.NewSchema(schemaConfig)
	if err == nil {
		l.schema = schema
	}

	return err
}

func (l *Logic) CreateMask(mask types.Mask) error {

	gqObject, err := converter.CreateObject(mask)
	if err != nil {
		return err
	}

	singleField, err := converter.MaskToSingleField(gqObject)
	if err != nil {
		return err
	}
	l.queryFields[mask.Name] = &singleField

	multipleField, err := converter.MaskToMultipleField(gqObject)
	if err != nil {
		return err
	}
	l.queryFields[fmt.Sprintf("%ss", mask.Name)] = &multipleField

	connectionField, err := converter.MaskToMultipleField(gqObject)
	if err != nil {
		return err
	}
	l.queryFields[fmt.Sprintf("%ssConnection", mask.Name)] = &connectionField

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

	queryFields := graphql.Fields{
		"_system": &graphql.Field{
			Type:        systemField,
			Description: "",
			Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
				return struct{}{}, nil
			},
		},
	}

	mutationFields := graphql.Fields{
		"createMask": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
				return nil, nil
			},
		},
	}

	query := graphql.ObjectConfig{Name: "Query", Fields: queryFields}
	mutation := graphql.ObjectConfig{Name: "Mutation", Fields: mutationFields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(query), Mutation: graphql.NewObject(mutation)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, err
	}

	return &Logic{
		queryFields:    queryFields,
		mutationFields: mutationFields,
		schema:         schema,
	}, nil
}
