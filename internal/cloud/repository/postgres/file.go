package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/baibikov/tensile-cloud/internal/cloud/types"
)

type File struct {
	db *sqlx.DB
}

func NewFile(db *sqlx.DB) *File {
	return &File{
		db: db,
	}
}

func (f File) SaveMeta(ctx context.Context, file types.File) (created types.File, err error) {
	query := `
		insert into files (folder_id, name, type, format)
		values (:folder_id, :name, :type, :format)
		returning id, type, name, format, created_at, updated_at
	`

	query, args, err := f.db.BindNamed(query, file)
	if err != nil {
		return types.File{}, err
	}

	return created, f.db.GetContext(ctx, &created, query, args...)
}

func (f File) RemoveMeta(ctx context.Context, id string) error {
	query := `
		delete from files where id = $1
	`

	_, err := f.db.ExecContext(ctx, query, id)
	return err
}

func (f File) GetByFolderID(ctx context.Context, folderID string, sort types.Sort) (files []types.File, err error) {
	query := `
		select 
		       id, 
		       type,
		       name,
		       format,
		       created_at,
		       updated_at 
		from files 
		where folder_id = $1 
		  and is_deleted = false
	`

	ssort := NewSort(sort, columnsDef{
		"name":      "name",
		"type":      "type",
		"createdAt": "created_at",
		"updatedAt": "updated_at",
	})

	query = ssort.OrderQuery(query)
	return files, f.db.SelectContext(ctx, &files, query, folderID)
}

func (f File) IsExists(ctx context.Context, id string) (ok bool, err error) {
	query := `
		select exists(
			select from files where id=$1
		)
	`

	return ok, f.db.GetContext(ctx, &ok, query, id)
}

func (f File) UpdateName(ctx context.Context, id, name string) (updated types.File, err error) {
	query := `
		update files set name=$2, updated_at=now() where id=$1
		returning id, type, name, format, created_at, updated_at
	`

	return updated, f.db.GetContext(ctx, &updated, query, id, name)
}

func (f File) UpdateFolderID(ctx context.Context, id []string, folderID string) error {
	query := `update files set folder_id=$1 where id=any($2)`

	_, err := f.db.ExecContext(ctx, query, folderID, pq.StringArray(id))
	return err
}

func (f File) GetByID(ctx context.Context, id string) (ff types.File, err error) {
	query := `select id, type, name, format, created_at, updated_at from files where id=$1`

	return ff, f.db.GetContext(ctx, &ff, query, id)
}

func (f File) MarkDelete(ctx context.Context, id []string) error {
	query := ` 
		update files set is_deleted=true, deleted_at = now() where id=any($2)
	`

	_, err := f.db.ExecContext(ctx, query, pq.StringArray(id))
	return err
}
