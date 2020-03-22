package converter

import (
	"github.com/bene/graphcms/logic/types"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"strings"
)

func convertType(from types.FieldType) *graphql.Scalar {

	switch from {
	case types.FieldTypeString:
		return graphql.String
	case types.FieldTypeInt:
		return graphql.Int
	case types.FieldTypeFloat:
		return graphql.Float
	case types.FieldTypeBool:
		return graphql.Boolean
	case types.FieldTypeDateTime:
		return graphql.DateTime
	case types.FieldTypeRichText:
		return graphql.String
	case types.FieldTypeMedia:
		return graphql.String
	}

	return nil
}

type selection struct {
	Name         string
	HasSubSelect bool
}

func reduceRowsOfOne(rows []map[string]interface{}) map[string]interface{} {

	result := make(map[string]interface{})
	for i, row := range rows {
		for column, value := range row {
			if strings.Contains(column, "___") {
				sub := strings.Split(column, "___")
				subTableName := sub[0]
				subCol := sub[1]

				if _, ok := result[subTableName]; ok {

					subEl := result[subTableName].([]map[string]interface{})
					if len(subEl) > i {
						subEl[i][subCol] = value
					} else {
						result[subTableName] = append(subEl, map[string]interface{}{
							subCol: value,
						})
					}

				} else {
					result[subTableName] = []map[string]interface{}{{
						subCol: value,
					}}
				}
			} else {
				result[column] = value
			}
		}
	}

	return result
}

func getSelectedFields(selectionPath []string, resolveParams graphql.ResolveParams) []selection {
	fields := resolveParams.Info.FieldASTs
	for _, propName := range selectionPath {
		found := false
		for _, field := range fields {
			if field.Name.Value == propName {
				selections := field.SelectionSet.Selections
				fields = make([]*ast.Field, 0)
				for _, selection := range selections {
					fields = append(fields, selection.(*ast.Field))
				}
				found = true
				break
			}
		}
		if !found {
			return []selection{}
		}
	}
	var collect []selection
	for _, field := range fields {
		collect = append(collect, selection{
			Name:         field.Name.Value,
			HasSubSelect: field.SelectionSet != nil,
		})
	}
	return collect
}
