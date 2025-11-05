package handlers

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed openapi-ui.html
var openapiui embed.FS

type openAPIUIHandler struct {
	content []byte
}

func NewOpenAPIUIHandler() (*openAPIUIHandler, error) {
	content, err := fs.ReadFile(openapiui, "openapi-ui.html")
	if err != nil {
		return nil, err
	}
	return &openAPIUIHandler{
		content: content,
	}, nil
}

func (h *openAPIUIHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(h.content)
}
