package cloud

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/baibikov/tensile-cloud/internal/cloud/types"
)

type FolderRepository interface {
	// Get getting folder by id
	Get(ctx context.Context, id string) (*types.Folder, error)
	// GetByParent recursive getting folders directory where parent id
	// if parent id is nil, then return root directory.
	GetByParent(ctx context.Context, parentID *string) ([]*types.Folder, error)
	// Create - create folder and return created folder
	Create(ctx context.Context, folder types.Folder) (types.Folder, error)
	// ExistsByID check exists folder where id return ok if folder exists
	// and return false if file is not exists
	ExistsByID(ctx context.Context, id string) (bool, error)
	// ExistsByParentIDName exists file/files by parent id and name
	// if parent id is nil, then check that name in root directory is exits
	ExistsByParentIDName(ctx context.Context, parentID *string, name string) (bool, error)
	// Update - update folder where id and return updated folder
	Update(ctx context.Context, folder types.Folder) (types.Folder, error)
	// Delete recursive remove folder and him directory folders like `bash: rm -r`
	Delete(ctx context.Context, id string) error
}

var (
	ErrFolderNotExists       = errors.New("folder not exists")
	ErrFolderParentNotExists = errors.New("parent not exists")
)

type Folder struct {
	folderrepo FolderRepository
}

func (f Folder) Find(ctx context.Context, parentID *string) ([]*types.Folder, error) {
	ff, err := f.folderrepo.GetByParent(ctx, parentID)
	return ff, errors.Wrap(err, "getting sub-folders")
}

func (f Folder) Folder(ctx context.Context, id string) (types.Folder, error) {
	ff, err := f.folderrepo.Get(ctx, id)
	if err != nil {
		return types.Folder{}, errors.Wrapf(err, "getting folder by id - %s", id)
	}
	if ff == nil {
		return types.Folder{}, ErrFolderNotExists
	}

	return *ff, nil
}

func (f Folder) Create(ctx context.Context, folder types.Folder) (types.Folder, error) {
	err := f.validParentID(ctx, folder.ParentID)
	if err != nil {
		return types.Folder{}, err
	}

	exists, err := f.folderrepo.ExistsByParentIDName(ctx, folder.ParentID, folder.Name)
	if err != nil {
		return types.Folder{}, errors.Wrap(err, "creating folder exists by parent/name")
	}
	if exists {
		return types.Folder{}, errors.Errorf("folder by name - %q has exists name is uniq state", folder.Name)
	}

	ff, err := f.folderrepo.Create(ctx, folder)
	if err != nil {
		return types.Folder{}, errors.Wrap(err, "create folder")
	}

	return ff, nil
}

func (f Folder) Update(ctx context.Context, folder types.Folder) (types.Folder, error) {
	err := f.validParentID(ctx, &folder.ID)
	if err != nil {
		return types.Folder{}, err
	}

	ff, err := f.folderrepo.Update(ctx, folder)
	if err != nil {
		return types.Folder{}, errors.Wrapf(err, "updating folder by id - %s", folder.ID)
	}

	return ff, nil
}

func (f Folder) validParentID(ctx context.Context, parentID *string) error {
	if parentID == nil {
		return nil
	}

	_, err := uuid.Parse(*parentID)
	if err != nil {
		return errors.Wrapf(err, "parse folder id with parent id - %s", *parentID)
	}

	return isFolderExists(ctx, f.folderrepo, *parentID)
}

func isFolderExists(ctx context.Context, folderrepo FolderRepository, id string) error {
	exists, err := folderrepo.ExistsByID(ctx, id)
	if err != nil {
		return errors.Wrap(err, "check existing parent id")
	}
	if !exists {
		return ErrFolderParentNotExists
	}

	return nil
}
