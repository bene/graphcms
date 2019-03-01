package types

type Field struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	Type        fieldType `json:"type"`
}
