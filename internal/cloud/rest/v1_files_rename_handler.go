package rest

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/models"
	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/ops/operations/files"
	"github.com/baibikov/tensile-cloud/internal/cloud/types"
	"github.com/baibikov/tensile-cloud/pkg/httperr"
)

func (h *Handler) RenameFileHandler(params files.RenameParams) middleware.Responder {
	updated, err := h.cloud.File.Rename(params.HTTPRequest.Context(), types.File{
		ID:   params.ID,
		Name: swag.StringValue(params.Body.Name),
	})
	if err != nil {
		return httperr.New().Bad(err)
	}

	return files.NewRenameOK().WithPayload(&models.File{
		ID:        updated.ID,
		Type:      updated.Type,
		Name:      updated.Name,
		Link:      updated.Link(),
		CreatedAt: updated.CreatedAt.Unix(),
		UpdatedAt: updated.UpdatedAt.Unix(),
	})
}
