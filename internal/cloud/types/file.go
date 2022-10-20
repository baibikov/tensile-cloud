package types

import (
	"io"
	"time"
)

type File struct {
	ID        string        `db:"id"`
	FolderID  string        `db:"folder_id"`
	Name      string        `db:"name"`
	Type      string        `db:"type"`
	Format    string        `db:"format"`
	CreatedAt time.Time     `db:"created_at"`
	UpdatedAt time.Time     `db:"updated_at"`
	Body      io.ReadCloser `db:"-"`
}

func (f File) FileName() string {
	return getFileName(f.Type, f.ID)
}

func (f File) Link() string {
	return filesUrl + f.FileName()
}

type Upload struct {
	ID        string
	Name      string
	Type      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u Upload) FileName() string {
	return getFileName(u.Type, u.ID)
}

func getFileName(typ, id string) string {
	if typ == "" {
		return id
	}

	return id + "." + typ
}

const filesUrl = "/files/"

func (u Upload) Link() string {
	return filesUrl + u.FileName()
}

type Move struct {
	FilesID  []string
	FolderID string
}

type DownloadFile struct {
	FileName string
	Payload  io.ReadCloser
}

type CopyFile struct {
	CopyID   string
	FolderID string
	Name     *string
}
