package gofip

import "fmt"

type ErrorAPI struct {
	Codigo      string `json:"code,omitempty"`
	Tipo        string `json:"type,omitempty"`
	Titulo      string `json:"title"`
	Estado      int    `json:"status"`
	Detalle     string `json:"detail"`
	Instancia   string `json:"instance,omitempty"`
	IDSolicitud string `json:"request_id,omitempty"`
}

func (e *ErrorAPI) Error() string {
	if e == nil {
		return ""
	}
	if e.Codigo != "" {
		return fmt.Sprintf("%s (%s): %s", e.Titulo, e.Codigo, e.Detalle)
	}
	if e.Titulo != "" {
		return fmt.Sprintf("%s: %s", e.Titulo, e.Detalle)
	}
	return e.Detalle
}

type ErrorHTTP struct {
	Estado int
	Cuerpo string
}

func (e *ErrorHTTP) Error() string {
	return fmt.Sprintf("respuesta HTTP %d: %s", e.Estado, e.Cuerpo)
}
