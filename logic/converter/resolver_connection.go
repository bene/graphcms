package converter

import (
	"github.com/graphql-go/graphql"
	_ "github.com/lib/pq"
)

func MaskToConnectionField(object *graphql.Object) (graphql.Field, error) {

	return graphql.Field{
		Type: graphql.NewList(object),
		Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
			return nil, nil
		},
	}, nil
}
