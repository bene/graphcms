package converter

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/bene/graphcms/logic/types"
	"github.com/graphql-go/graphql"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

func MaskToSingleField(mask types.Mask) (graphql.Field, error) {

	db, err := sql.Open("postgres", "postgres://postgres:mysecretpassword@localhost:5431/evobend?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

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
					Type:        graphql.NewNonNull(graphql.NewList(convertType(repeatingType))),
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
			Type:        convertType(field.Type),
			Description: field.Description,
		}
	}

	return graphql.Field{
		Args:        args,
		Type:        maskObject,
		Description: fmt.Sprintf("Query a single %s. %s", mask.DisplayName, mask.Description),
		Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {

			tableName := p.Info.FieldName
			tableFields := []string{fmt.Sprintf("%s.id", tableName)}
			tableJoins := []string{}

			selections := getSelectedFields([]string{mask.Name}, p)
			for _, s := range selections {
				if !s.HasSubSelect {
					tableFields = append(tableFields, fmt.Sprintf("%s.%s", tableName, s.Name))
				} else {

					// TODO: Make recursive for deeper levels
					subSelections := getSelectedFields([]string{mask.Name, s.Name}, p)
					for _, ss := range subSelections {
						tableFields = append(tableFields, fmt.Sprintf("%s.%s AS %s___%s", s.Name, ss.Name, s.Name, ss.Name))
					}
					tableJoins = append(tableJoins, fmt.Sprintf(" LEFT JOIN %s_%s %s ON %s.id = %s.%s_id", tableName, s.Name, s.Name, tableName, s.Name, tableName))
				}
			}

			query := fmt.Sprintf("SELECT %s FROM %s%s", strings.Join(tableFields, ", "), tableName, strings.Join(tableJoins, " "))
			log.Println(query)

			rows, err := db.QueryContext(p.Context, query)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			cols, err := rows.Columns()
			if err != nil {
				return nil, err
			}

			for _, col := range cols {
				fmt.Printf("Col: %s", col)
			}

			results := map[string][]map[string]interface{}{}
			for rows.Next() {

				columns := make([]interface{}, len(cols))
				columnPointers := make([]interface{}, len(cols))
				for i, _ := range columns {
					columnPointers[i] = &columns[i]
				}

				// Scan the result into the column pointers...
				if err := rows.Scan(columnPointers...); err != nil {
					return nil, err
				}

				m := make(map[string]interface{})
				for i, colName := range cols {
					val := columnPointers[i].(*interface{})
					m[colName] = *val
				}

				id := fmt.Sprint(m["id"])

				if _, ok := results[id]; ok {
					results[id] = append(results[id], m)
				} else {
					results[id] = []map[string]interface{}{m}
				}
			}

			reduceds := []map[string]interface{}{}

			for _, one := range results {
				reduceds = append(reduceds, reduceRowsOfOne(one))
			}

			return reduceds[0], nil
		},
	}, nil
}
