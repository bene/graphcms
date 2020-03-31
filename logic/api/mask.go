package api

import "github.com/graphql-go/graphql"

var MaskType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Mask",
	Description: "",
	Interfaces:  nil,
	Fields:      nil,
	IsTypeOf:    nil,
})

var apiMutationFields = graphql.Fields{
	"createMask": &graphql.Field{
		Type: MaskType,
		Args: nil,
		Resolve: func(p graphql.ResolveParams) (i interface{}, err error) {
			return nil, nil
		},
		DeprecationReason: "",
		Description:       "",
	},
}
