package types

type FieldType string

func (f FieldType) GetTypeName() string {
	return string(f)
}

var (
	FieldTypeString    FieldType = "STRING"
	FieldTypeInt       FieldType = "INT"
	FieldTypeFloat     FieldType = "FLOAT"
	FieldTypeBool      FieldType = "BOOL"
	FieldTypeDateTime  FieldType = "DATETIME"
	FieldTypeRichText  FieldType = "RICH_TEXT"
	FieldTypeMedia     FieldType = "MEDIA"
	FieldTypeRepeating FieldType = "REPEATING"
)
