package types

type Mask struct {
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	Description string  `json:"description"`
	Fields      []Field `json:"fields"`
}
