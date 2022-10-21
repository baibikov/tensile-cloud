package rest

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/ops/operations/files"
	"github.com/baibikov/tensile-cloud/pkg/httperr"
)

func (h *Handler) MarkDeleteFileHandler(params files.MarkDeleteParams) middleware.Responder {
	err := h.cloud.File.MarkDelete(params.HTTPRequest.Context(), params.Body.FilesID)
	if err != nil {
		return httperr.New().Bad(err)
	}

	return files.NewMarkDeleteOK()
}
