package converter

import (
	"errors"
	"fmt"
	"github.com/bene/graphcms/logic/types"
	"github.com/graphql-go/graphql"
	"log"
)

func MaskToSingleField(mask types.Mask) (graphql.Field, error) {

	fields := graphql.Fields{}
	for _, field := range mask.Fields {

		if field.Type == types.FieldTypeRepeating {
			if repeatingMask, ok := field.Of.(types.Mask); ok {

				subField, err := MaskToSingleField(repeatingMask)
				if err != nil {
					return graphql.Field{}, err
				}

				fields[field.Name] = &graphql.Field{
					Type:        graphql.NewNonNull(graphql.NewList(subField.Type)),
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

			} else if repeatingType, ok := field.Of.(types.FieldType); ok {

				fields[field.Name] = &graphql.Field{
					Type:        graphql.NewNonNull(graphql.NewList(ConvertType(repeatingType))),
					Description: field.Description,
					Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
						if result, ok := p.Source.(map[string]interface{}); ok {
							if value, ok := result[field.Name]; ok {
								return value, nil
							}
						}
						return nil, nil
					},
				}
			} else {
				return graphql.Field{}, errors.New("invalid repeating")
			}
		} else {

			var t graphql.Output
			if field.IsRequired {
				t = graphql.NewNonNull(ConvertType(field.Type))
			} else {
				t = ConvertType(field.Type)
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

	maskObject := graphql.NewObject(graphql.ObjectConfig{
		Name:   mask.DisplayName,
		Fields: fields,
	})

	// Map unique fields as arguments
	args := make(map[string]*graphql.ArgumentConfig)
	for _, field := range mask.Fields {
		if !field.IsRequired || field.Type == types.FieldTypeRepeating {
			continue
		}

		args[field.Name] = &graphql.ArgumentConfig{
			Type:        ConvertType(field.Type),
			Description: field.Description,
		}
	}

	return graphql.Field{
		Args:        args,
		Type:        maskObject,
		Description: fmt.Sprintf("Query a single %s. %s", mask.DisplayName, mask.Description),
		Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
			// TODO: Get from database

			log.Println("Request")
			log.Print(len(p.Info.FieldASTs))

			return map[string]interface{}{
				"name": "Unternehmen",
				"no_index": 0,
				"menu_index": "false",
				"content": []interface{}{
					map[string]interface{}{
						"name": "title",
						"value": "Unternehmen",
					},
					map[string]interface{}{
						"name": "description",
						"value": "Über uns",
					},
				},
			}, nil
		},
	}, nil
}
