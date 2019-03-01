package types

type fieldType string

func (f fieldType) GetTypeName() string {
	return string(f)
}

var (
	FieldTypeString   fieldType = "string"
	FieldTypeInt      fieldType = "int"
	FieldTypeDate     fieldType = "date"
	FieldTypeDateBool fieldType = "bool"
	FieldTypeRichText fieldType = "rich_text"
	FieldTypeMedia    fieldType = "media"
	FieldTypeTimeDate fieldType = "time_date"
)
