package models

import (
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID       uuid.UUID `bun:"id,type:uuid" json:"id"`
	Username string    `bun:"username" json:"username"`
	Password string    `bun:"password" json:"-"`
	Email    string    `bun:"email" json:"email"`
}

func (u *User) GetID() string {
	return u.ID.String()
}
