package rest

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"

	"github.com/baibikov/tensile-cloud/internal/cloud/rest/generated/ops/operations/files"
	"github.com/baibikov/tensile-cloud/pkg/httperr"
)

func (h *Handler) V1DownloadFileHandler(params files.DownloadParams) middleware.Responder {
	res, err := h.cloud.File.Download(params.HTTPRequest.Context(), params.ID)
	if err != nil {
		return httperr.New().Bad(err)
	}

	api := files.NewDownloadOK()
	api.SetContentDisposition(fmt.Sprintf("attachment; filename=%q", res.FileName))
	return api.WithPayload(res.Payload)
}
