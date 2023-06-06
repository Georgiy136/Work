package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Project struct {
	bun.BaseModel `bun:"table:projects"`

	Id          uuid.UUID   `bun:"uuid"`
	ProjectName string      `bun:"project_name"`
	ProjectType string      `bun:"project_type"`
	Operators   []*Operator `bun:"-"`
}
