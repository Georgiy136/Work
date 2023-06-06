package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Client struct {
	bun.BaseModel `bun:"table:clients"`

	Id         uuid.UUID `bun:"uuid"`
	ClientRole string    `bun:"role"`
	Username   string    `bun:"username"`
	Login      string    `bun:"login"`
	Password   string    `bun:"password"`
}
