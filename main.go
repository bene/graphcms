package main

import (
	"github.com/bene/graphcms/logic"
	"github.com/bene/graphcms/logic/types"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"log"
)

func main() {
	e := echo.New()

	e.GET("/", func(context echo.Context) error {

		project, err := logic.CreateProject("example-1", uuid.New())
		if err != nil {
			return err
		}

		project.Models = []types.Model{
			{
				Name:        "article",
				DisplayName: "Article",
				Fields: []types.Field{
					{
						Name:        "cover_image",
						DisplayName: "Cover Image",
						Type:        types.FieldTypeMedia,
					},
					{
						Name:        "title",
						DisplayName: "Title",
						Type:        types.FieldTypeString,
					},
					{
						Name:        "author",
						DisplayName: "Author",
						Type:        types.FieldTypeString,
					},
					{
						Name:        "content",
						DisplayName: "Content",
						Type:        types.FieldTypeString,
					},
				},
			},
			{
				Name:        "event",
				DisplayName: "Event",
				Fields: []types.Field{
					{
						Name:        "title",
						DisplayName: "Title",
						Type:        types.FieldTypeString,
					},
					{
						Name:        "start",
						DisplayName: "Start",
						Type:        types.FieldTypeTimeDate,
					},
					{
						Name:        "end",
						DisplayName: "End",
						Type:        types.FieldTypeTimeDate,
					},
					{
						Name:        "location",
						DisplayName: "Location",
						Type:        types.FieldTypeTimeDate,
					},
				},
			},
		}

		return context.JSON(200, project)
	})

	log.Fatalln(e.Start(":2000"))
}
