package rest

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/models"
	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/ops/operations/files"
	"github.com/baibikov/tensile-cloud/internal/cloud/types"
	"github.com/baibikov/tensile-cloud/pkg/httperr"
)

func (h *Handler) V1CopyFileHandler(params files.CopyParams) middleware.Responder {
	file, err := h.cloud.File.Copy(params.HTTPRequest.Context(), types.CopyFile{
		CopyID:   swag.StringValue(params.Body.CopyID),
		FolderID: swag.StringValue(params.Body.FolderID),
		Name:     swag.String(params.Body.Name),
	})
	if err != nil {
		return httperr.New().Bad(err)
	}

	return files.NewCopyOK().WithPayload(&models.File{
		ID:        file.ID,
		Type:      file.Type,
		Name:      file.Name,
		Link:      file.Link(),
		CreatedAt: file.CreatedAt.Unix(),
		UpdatedAt: file.UpdatedAt.Unix(),
	})
}
