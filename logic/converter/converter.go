package converter

import (
	"errors"
	"github.com/bene/graphcms/logic/types"
	"github.com/graphql-go/graphql"
)

func createListField(repeating types.Field) (*graphql.Field, error) {

	if repeatingMask, ok := repeating.Of.(types.Mask); ok {

		object, err := CreateObject(repeatingMask)
		if err != nil {
			return nil, err
		}

		return &graphql.Field{
			Name:        repeating.Name,
			Type:        graphql.NewNonNull(graphql.NewList(object)),
			Description: repeating.Description,
			Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
				if result, ok := p.Source.(map[string]interface{}); ok {

					if value, ok := result[p.Info.FieldName]; ok {
						return value, nil
					}
				}
				return nil, nil
			},
		}, nil

	} else if fieldType, ok := repeating.Of.(types.FieldType); ok {

		return &graphql.Field{
			Name:        repeating.Name,
			Type:        graphql.NewNonNull(graphql.NewList(convertType(fieldType))),
			Description: repeating.Description,
			Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
				if result, ok := p.Source.(map[string]interface{}); ok {

					if value, ok := result[p.Info.FieldName]; ok {
						return value, nil
					}
				}
				return nil, nil
			},
		}, nil

	} else {
		return nil, errors.New("invalid of type")
	}
}

func CreateObject(mask types.Mask) (*graphql.Object, error) {

	fields := graphql.Fields{}
	for _, field := range mask.Fields {

		if field.Type == types.FieldTypeRepeating {

			repeatingField, err := createListField(field)
			if err != nil {
				return nil, err
			}
			fields[field.Name] = repeatingField

		} else {

			var t graphql.Output
			if field.IsRequired {
				t = graphql.NewNonNull(convertType(field.Type))
			} else {
				t = convertType(field.Type)
			}

			fields[field.Name] = &graphql.Field{
				Name:        field.Name,
				Type:        t,
				Description: field.Description,
				Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
					if result, ok := p.Source.(map[string]interface{}); ok {

						if value, ok := result[p.Info.FieldName]; ok {
							return value, nil
						}
					}
					return nil, nil
				},
			}
		}
	}

	return graphql.NewObject(graphql.ObjectConfig{
		Name:        mask.DisplayName,
		Description: mask.Description,
		Fields:      fields,
	}), nil
}
