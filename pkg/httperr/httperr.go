package httperr

import (
	"encoding/json"
	"net/http"

	"github.com/go-openapi/runtime"
)

type Err struct {
}

type Responder struct {
	e    error
	code int32
}

type message struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func (r *Responder) WriteResponse(w http.ResponseWriter, _ runtime.Producer) {
	bb, err := json.Marshal(message{
		Code:    r.code,
		Message: r.e.Error(),
	})
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	_, _ = w.Write(bb)
}

func (r Err) Bad(err error) *Responder {
	return &Responder{
		e:    err,
		code: http.StatusBadRequest,
	}
}

func New() *Err {
	return &Err{}
}
