package rest

import (
	"context"
	"net/http"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"

	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/ops"
	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/ops/operations"
	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/ops/operations/directory"
	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/ops/operations/files"
	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/ops/operations/folder"
	"github.com/baibikov/tensile-cloud/internal/cloud/types"
)

type Folder interface {
	Folder(ctx context.Context, id string) (types.Folder, error)
	Find(ctx context.Context, parentID *string) ([]*types.Folder, error)
	Create(ctx context.Context, folder types.Folder) (types.Folder, error)
	Update(ctx context.Context, folder types.Folder) (types.Folder, error)
}

type File interface {
	Upload(ctx context.Context, files []types.File) ([]types.Upload, error)
	Find(ctx context.Context, folderID string) ([]types.File, error)
	Rename(ctx context.Context, file types.File) (types.File, error)
	Move(ctx context.Context, move types.Move) error
	Download(ctx context.Context, id string) (types.DownloadFile, error)
	Copy(ctx context.Context, file types.CopyFile) (created types.File, err error)
}

type Config struct {
	UploadSize int64
}

type UseCase struct {
	Folder Folder
	File   File
}

type Handler struct {
	api    *operations.ClouderAPI
	cloud  *UseCase
	config *Config
}

func (h *Handler) Handler(builder middleware.Builder) http.Handler {
	return h.api.Serve(builder)
}

func New(config *Config, cloud *UseCase) (*Handler, error) {
	spec, err := loads.Analyzed(ops.SwaggerJSON, "")
	if err != nil {
		return nil, err
	}

	h := &Handler{
		api:    operations.NewClouderAPI(spec),
		cloud:  cloud,
		config: config,
	}

	// folder
	h.api.FolderCreateHandler = folder.CreateHandlerFunc(h.FolderCreateHandler)
	h.api.FolderUpdateHandler = folder.UpdateHandlerFunc(h.FolderUpdateHandler)

	// directory
	h.api.DirectoryGetDirectoryHandler = directory.GetDirectoryHandlerFunc(h.DirectoryListHandler)

	// files
	h.api.FilesUploadHandler = files.UploadHandlerFunc(h.UploadFilesHandler)
	h.api.FilesRenameHandler = files.RenameHandlerFunc(h.RenameFileHandler)
	h.api.FilesMoveHandler = files.MoveHandlerFunc(h.MoveFileHandler)
	h.api.FilesDownloadHandler = files.DownloadHandlerFunc(h.DownloadFileHandler)
	h.api.FilesCopyHandler = files.CopyHandlerFunc(h.CopyFileHandler)

	return h, nil
}
