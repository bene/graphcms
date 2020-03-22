package main

import (
	"github.com/bene/graphcms/logic"
	"github.com/bene/graphcms/logic/types"
	"github.com/graphql-go/handler"
	"log"
	"net/http"
)

func main() {

	l, err := logic.NewLogic()
	if err != nil {
		log.Panic(err)
	}

	h := handler.New(&handler.Config{
		Schema:   l.GetSchema(),
		Pretty:   true,
		GraphiQL: true,
	})

	err = l.CreateMask(types.Mask{
		Name:        "page",
		DisplayName: "Page",
		Description: "Pages are the main components of an website.",
		Fields: []types.Field{
			{Name: "slug", DisplayName: "Slug", Description: "URL of the page.", UseAsTitle: false, IsUnique: true, IsRequired: true, Localize: true, Type: types.FieldTypeString},
			{Name: "title", DisplayName: "Title", Description: "Title, important for SEO.", UseAsTitle: true, IsUnique: false, IsRequired: true, Localize: true, Type: types.FieldTypeString},
			{Name: "content", DisplayName: "Content", Description: "Page content", UseAsTitle: false, IsUnique: false, IsRequired: true, Localize: true, Type: types.FieldTypeRepeating, Of: types.Mask{
				Name:        "content",
				DisplayName: "Content",
				Description: "Content data",
				Fields: []types.Field{
					{Name: "name", DisplayName: "Name", Description: "Name of the data", UseAsTitle: true, IsUnique: false, IsRequired: true, Localize: false, Type: types.FieldTypeString},
					{Name: "value", DisplayName: "Value", Description: "Value of the data", UseAsTitle: false, IsUnique: false, IsRequired: true, Localize: false, Type: types.FieldTypeString},
				},
			}},
			{Name: "tags", DisplayName: "Tags", Description: "Page tags", UseAsTitle: false, IsUnique: false, IsRequired: false, Localize: true, Type: types.FieldTypeRepeating, Of: types.FieldTypeString},
		},
	})
	if err != nil {
		log.Println(err)
	}

	http.Handle("/graph", h)
	err = http.ListenAndServe(":4242", nil)
	if err != nil {
		log.Panic(err)
	}
}
