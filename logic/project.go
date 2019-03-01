package logic

import (
	"github.com/bene/graphcms/logic/types"
	"github.com/google/uuid"
	"time"
)

func CreateProject(name string, owner uuid.UUID) (*types.ProjectMetadata, error) {

	id := uuid.New()
	now := time.Now()

	project := types.ProjectMetadata{
		ID:        id,
		Name:      name,
		CreatedAt: now,
		Owner:     owner,
		Models:    []types.Model{},
	}

	// TODO: Save in db

	return &project, nil
}
