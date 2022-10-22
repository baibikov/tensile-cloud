package rest

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/models"
	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/ops/operations/folder"
	"github.com/baibikov/tensile-cloud/internal/cloud/types"
	"github.com/baibikov/tensile-cloud/pkg/httperr"
)

func (h *Handler) V1FolderUpdateHandler(params folder.UpdateParams) middleware.Responder {
	res, err := h.cloud.Folder.Update(params.HTTPRequest.Context(), types.Folder{
		ID:   params.ID,
		Name: swag.StringValue(params.Body.Name),
	})
	if err != nil {
		return httperr.New().Bad(err)
	}

	return folder.NewCreateOK().WithPayload(&models.Folder{
		ID:        res.ID,
		Name:      res.Name,
		ParentID:  res.ParentID,
		CreatedAt: res.CreatedAt.Unix(),
		UpdatedAt: res.UpdatedAt.Unix(),
	})
}
