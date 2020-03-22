package converter

import (
	"github.com/bene/graphcms/logic/types"
	"github.com/graphql-go/graphql"
)

func ConvertType(from types.FieldType) *graphql.Scalar {

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
