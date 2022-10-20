package rest

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/models"
	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/ops/operations/directory"
	"github.com/baibikov/tensile-cloud/pkg/httperr"
)

func (h *Handler) DirectoryListHandler(params directory.GetDirectoryParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()
	resfolders, err := h.cloud.Folder.Find(ctx, params.ParentID)
	if err != nil {
		return httperr.New().Bad(err)
	}

	dir := make([]*models.Folder, 0, len(resfolders))
	for _, r := range resfolders {
		dir = append(dir, &models.Folder{
			ID:        r.ID,
			Name:      r.Name,
			ParentID:  r.ParentID,
			CreatedAt: r.CreatedAt.Unix(),
			UpdatedAt: r.UpdatedAt.Unix(),
		})
	}

	if params.ParentID == nil {
		return directory.NewGetDirectoryOK().WithPayload(&models.Directory{
			Folders: dir,
		})
	}

	resfiles, err := h.cloud.File.Find(ctx, *params.ParentID)
	if err != nil {
		return httperr.New().Bad(err)
	}

	ff := make([]*models.File, 0, len(resfiles))
	for _, f := range resfiles {
		ff = append(ff, &models.File{
			ID:        f.ID,
			Type:      f.Type,
			Name:      f.Name,
			Link:      f.Link(),
			CreatedAt: f.CreatedAt.Unix(),
			UpdatedAt: f.UpdatedAt.Unix(),
		})
	}

	return directory.NewGetDirectoryOK().WithPayload(&models.Directory{
		Folders: dir,
		Files:   ff,
	})
}
