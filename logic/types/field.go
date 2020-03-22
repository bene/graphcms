package types

type Field struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	UseAsTitle  bool      `json:"use_as_title"`
	IsUnique    bool      `json:"is_unique"`
	IsRequired  bool      `json:"is_required"`
	Localize    bool      `json:"localize"`
	Type        FieldType `json:"type"`

	// Defines the type (must be either instance of mask or fieldType
	Of interface{} `json:"of"`

	// Reference of relation
	Ref string
}
