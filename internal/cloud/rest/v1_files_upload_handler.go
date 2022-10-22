package rest

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/pkg/errors"

	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/models"
	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/ops/operations/files"
	"github.com/baibikov/tensile-cloud/internal/cloud/types"
	"github.com/baibikov/tensile-cloud/pkg/httperr"
)

const (
	filesMultiPartFormKey = "files"
)

func (h *Handler) V1UploadFilesHandler(params files.UploadParams) middleware.Responder {
	err := params.HTTPRequest.ParseMultipartForm(h.config.UploadSize)
	if err != nil {
		return httperr.New().Bad(err)
	}

	uploadFiles, ok := params.HTTPRequest.MultipartForm.File[filesMultiPartFormKey]
	if !ok {
		return httperr.New().Bad(errors.New("param files is required"))
	}

	ff := make([]types.File, 0, len(uploadFiles))

	for _, f := range uploadFiles {
		fo, err := f.Open()
		if err != nil {
			return httperr.New().Bad(err)
		}

		ff = append(ff, types.File{
			FolderID: params.FolderID,
			Name:     f.Filename,
			Body:     fo,
		})
	}

	created, err := h.cloud.File.Upload(params.HTTPRequest.Context(), ff)
	if err != nil {
		return httperr.New().Bad(err)
	}

	res := make([]*models.File, 0, len(created))
	for _, f := range created {
		res = append(res, &models.File{
			ID:        f.ID,
			Type:      f.Type,
			Name:      f.Name,
			Link:      f.Link(),
			CreatedAt: f.CreatedAt.Unix(),
			UpdatedAt: f.UpdatedAt.Unix(),
		})
	}

	return files.NewUploadOK().WithPayload(&models.CreatedFiles{
		Files: res,
	})
}
