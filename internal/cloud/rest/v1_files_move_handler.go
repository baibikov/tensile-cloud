package rest

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"

	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/ops/operations/files"
	"github.com/baibikov/tensile-cloud/internal/cloud/types"
	"github.com/baibikov/tensile-cloud/pkg/httperr"
)

func (h *Handler) MoveFileHandler(params files.MoveParams) middleware.Responder {
	err := h.cloud.File.Move(params.HTTPRequest.Context(), types.Move{
		FilesID:  params.Body.FilesID,
		FolderID: swag.StringValue(params.Body.FolderID),
	})
	if err != nil {
		return httperr.New().Bad(err)
	}

	return files.NewMoveOK()
}
