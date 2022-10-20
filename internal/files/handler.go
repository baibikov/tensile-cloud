package files

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

type Handler struct {
	dir string
}

func New(dir string) *Handler {
	return &Handler{
		dir: dir,
	}
}

func (h *Handler) Handler(_ middleware.Builder) http.Handler {
	return http.StripPrefix("/files/", http.FileServer(http.Dir(h.dir)))
}
