package types

import (
	"github.com/google/uuid"
	"time"
)

type ProjectMetadata struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Owner     uuid.UUID `json:"owner"`
	Models    []Model   `json:"models"`
}

type Model struct {
	Name        string  `json:"name"`
	DisplayName string  `json:"display_name"`
	Fields      []Field `json:"fields"`
}
