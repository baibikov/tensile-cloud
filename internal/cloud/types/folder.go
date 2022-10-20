package types

import (
	"time"
)

type Folder struct {
	ID        string    `db:"id"`
	ParentID  *string   `db:"parent_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
