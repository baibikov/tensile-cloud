package cloud

import (
	"context"
	"io"

	"github.com/pkg/errors"
	"go.uber.org/multierr"

	"github.com/baibikov/tensile-cloud/internal/cloud/types"
	"github.com/baibikov/tensile-cloud/pkg/utils"
)

type FileRepository interface {
	// IsExists check file exists by id (file id)
	IsExists(ctx context.Context, id string) (ok bool, err error)
	// GetByFolderID return files info by folderID (folder binding)
	GetByFolderID(ctx context.Context, folderID string, sort types.Sort) ([]types.File, error)
	// UpdateName - update file name where file id
	// return updated file info
	UpdateName(ctx context.Context, id, name string) (updated types.File, err error)
	// UpdateFolderID - update files folder id by files id
	UpdateFolderID(ctx context.Context, id []string, folderID string) error
	// Open opening file by name with full file information like: file.txt,
	// without full dir path
	// use with file information id like: 4d2e16db-3bb4-430d-b7e6-6844da203595.pdf
	Open(ctx context.Context, name string) (io.ReadCloser, error)
	// Copy copied some file
	// dst is a destination file name with file information (with type)
	// src is a source file name with file information (with type)
	Copy(ctx context.Context, dst, src string) (err error)
	// GetByID getting file by id
	GetByID(ctx context.Context, id string) (ff types.File, err error)
	// MarkDelete is a safety delete file info
	// don't delete some file by mark delete
	MarkDelete(ctx context.Context, id []string) error
}

var (
	ErrFileNotExists = errors.New("file not exists")
)

type File struct {
	saver      Saver
	filerepo   FileRepository
	folderrepo FolderRepository
}

func (f File) Upload(ctx context.Context, files []types.File) ([]types.Upload, error) {
	if len(files) == 0 {
		return nil, nil
	}

	err := isFolderExists(ctx, f.folderrepo, files[0].FolderID)
	if err != nil {
		return nil, err
	}

	res := make([]types.Upload, 0, len(files))
	for _, file := range files {
		created, err := f.uploadSave(ctx, file)
		if err != nil {
			return nil, err
		}

		res = append(res, created)
	}

	return res, nil
}

func (f File) uploadSave(ctx context.Context, file types.File) (upload types.Upload, err error) {
	defer func() {
		if err != nil {
			multierr.AppendInto(&err, f.rollbackSavedFile(ctx, upload.ID, upload.FileName()))
		}

		multierr.AppendInto(&err, file.Body.Close())
	}()

	file.Type = utils.FileType(file.Name)

	created, err := f.saver.SaveMeta(ctx, file)
	if err != nil {
		return types.Upload{}, errors.Wrapf(err, "save file meta by name - %s", file.Name)
	}

	err = f.saver.SaveFile(ctx, types.File{
		ID:   created.ID,
		Name: created.Name,
		Type: created.Type,
		Body: file.Body,
	})
	if err != nil {
		return types.Upload{}, errors.Wrapf(err, "save file body by name - %s", file.Name)
	}

	return types.Upload{
		ID:        created.ID,
		Name:      created.Name,
		Type:      created.Type,
		CreatedAt: created.CreatedAt,
		UpdatedAt: created.UpdatedAt,
	}, nil
}

func (f File) rollbackSavedFile(ctx context.Context, id, filename string) error {
	if id == "" {
		return nil
	}

	var (
		resErr error
	)

	err := f.saver.RemoveMeta(ctx, id)
	if err != nil {
		resErr = multierr.Append(resErr, err)
	}

	err = f.saver.RemoveFile(ctx, filename)
	if err != nil {
		resErr = multierr.Append(resErr, err)
	}

	return resErr
}

func (f File) Find(ctx context.Context, folderID string, sort types.Sort) ([]types.File, error) {
	err := isFolderExists(ctx, f.folderrepo, folderID)
	if err != nil {
		return nil, err
	}

	ff, err := f.filerepo.GetByFolderID(ctx, folderID, sort)
	return ff, errors.Wrapf(err, "getting files by folder id - %s", folderID)
}

func (f File) Rename(ctx context.Context, file types.File) (types.File, error) {
	err := isFileExists(ctx, f.filerepo, file.ID)
	if err != nil {
		return types.File{}, err
	}

	ff, err := f.filerepo.UpdateName(ctx, file.ID, file.Name)
	return ff, errors.Wrapf(err, "renaming file by id - %s, name - %s", file.ID, file.Name)
}

func (f File) Move(ctx context.Context, move types.Move) error {
	err := isFolderExists(ctx, f.folderrepo, move.FolderID)
	if err != nil {
		return err
	}

	return errors.Wrapf(
		f.filerepo.UpdateFolderID(ctx, move.FilesID, move.FolderID),
		"moving files by folder id - %s",
		move.FolderID,
	)
}

func (f File) Download(ctx context.Context, id string) (types.DownloadFile, error) {
	err := isFileExists(ctx, f.filerepo, id)
	if err != nil {
		return types.DownloadFile{}, err
	}

	file, err := f.filerepo.GetByID(ctx, id)
	if err != nil {
		return types.DownloadFile{}, errors.Wrapf(err, "getting file by id - %s", id)
	}

	iio, err := f.filerepo.Open(ctx, file.FileName())
	if err != nil {
		return types.DownloadFile{}, errors.Wrapf(err, "downloading file by id - %s", id)
	}

	return types.DownloadFile{
		FileName: file.Name,
		Payload:  iio,
	}, nil
}

func (f File) Copy(ctx context.Context, fcopy types.CopyFile) (created types.File, err error) {
	err = isFileExists(ctx, f.filerepo, fcopy.CopyID)
	if err != nil {
		return types.File{}, err
	}
	defer func() {
		if err == nil {
			return
		}

		multierr.AppendInto(&err, f.rollbackSavedFile(ctx, created.ID, created.FileName()))
	}()

	src, err := f.filerepo.GetByID(ctx, fcopy.CopyID)
	if err != nil {
		return types.File{}, err
	}

	name := src.Name
	if fcopy.Name != nil {
		name = *fcopy.Name
	}

	created, err = f.saver.SaveMeta(ctx, types.File{
		FolderID: fcopy.FolderID,
		Name:     name,
		Type:     src.Type,
	})
	if err != nil {
		return types.File{}, err
	}

	return created, f.filerepo.Copy(ctx, created.FileName(), src.FileName())
}

func (f File) MarkDelete(ctx context.Context, id []string) error {
	return errors.Wrap(
		f.filerepo.MarkDelete(ctx, id),
		"deleted files",
	)
}

func isFileExists(ctx context.Context, repository FileRepository, id string) error {
	ok, err := repository.IsExists(ctx, id)
	if err != nil {
		return err
	}
	if !ok {
		return ErrFileNotExists
	}

	return nil
}
