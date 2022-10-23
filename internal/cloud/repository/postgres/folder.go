package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"

	"github.com/baibikov/tensile-cloud/internal/cloud/types"
)

type Folder struct {
	db *sqlx.DB
	sb *sqlbuilder.SelectBuilder
}

func NewFolder(db *sqlx.DB) *Folder {
	return &Folder{
		db: db,
		sb: sqlbuilder.PostgreSQL.NewSelectBuilder(),
	}
}

func (f Folder) ExistsByID(ctx context.Context, parentID string) (exists bool, err error) {
	query := `
		select exists(
			select from folders where id=$1
		)
	`

	return exists, f.db.GetContext(ctx, &exists, query, parentID)
}

func (f Folder) ExistsByParentIDName(ctx context.Context, parentID *string, name string) (exists bool, err error) {
	if parentID == nil {
		return exists, f.db.GetContext(
			ctx,
			&exists,
			"select exists(select from folders where parent_id is null and name=$1)",
			name,
		)
	}

	return exists, f.db.GetContext(
		ctx,
		&exists,
		"select exists(select from folders where parent_id = $1 and name=$2)",
		parentID,
		name,
	)
}

func (f Folder) Get(ctx context.Context, id string) (*types.Folder, error) {
	query := `
		 select
			   id,
			   parent_id,
			   name,
			   created_at,
			   updated_at
		from folders
		where id = $1
	`

	folder := &types.Folder{}
	err := f.db.GetContext(ctx, folder, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return folder, err
}

func (f Folder) GetByParent(ctx context.Context, parentID *string, sort types.Sort) (folders []*types.Folder, err error) {
	sb := f.sb.
		Select(
			"id",
			"parent_id",
			"name",
			"created_at",
			"updated_at",
		).
		From("folders")

	if parentID != nil {
		sb = sb.Where(sb.IsNull("parent_id"))
	} else {
		sb = sb.Where(sb.Equal("parent_id", *parentID))
	}
	query, args := sb.Build()
	query = NewSort(sort, newColumnsDef("name", "created_at", "updated_at")).OrderQuery(query)

	return folders, f.db.SelectContext(ctx, &folders, query, args...)
}

func (f Folder) Create(ctx context.Context, folder types.Folder) (ff types.Folder, err error) {
	query := `
		insert into folders(name, parent_id)
		values ($1, $2)
		returning id, parent_id, name, created_at, updated_at
	`

	return ff, f.db.GetContext(ctx, &ff, query, folder.Name, folder.ParentID)
}

func (f Folder) Update(ctx context.Context, folder types.Folder) (ff types.Folder, err error) {
	query := `
		update folders set name = $1, updated_at = now() where id= $2
		returning id, parent_id, name, created_at, updated_at
	`

	return ff, f.db.GetContext(ctx, &ff, query, folder.Name, folder.ID)
}

func (f Folder) Delete(ctx context.Context, id string) error {
	panic(any("implement me"))
}
