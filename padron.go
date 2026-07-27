package gofip

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/padron"
)

type ServicioPadron struct {
	transporte *transporte
}

func (s *ServicioPadron) ConsultarP4(ctx context.Context, cuit string) (*padron.RespuestaPersonaPadron, error) {
	return s.consultar(ctx, "/api/v2/arca/padron/p4/personas/"+url.PathEscape(cuit))
}

func (s *ServicioPadron) ConsultarP5(ctx context.Context, cuit string) (*padron.RespuestaPersonaPadron, error) {
	return s.consultar(ctx, "/api/v2/arca/padron/p5/personas/"+url.PathEscape(cuit))
}

func (s *ServicioPadron) ConsultarP10(ctx context.Context, cuit string) (*padron.RespuestaPersonaPadron, error) {
	return s.consultar(ctx, "/api/v2/arca/padron/p10/personas/"+url.PathEscape(cuit))
}

func (s *ServicioPadron) ConsultarP13(
	ctx context.Context,
	tipoDocumento padron.TipoDocumentoP13,
	identificador string,
) (*padron.RespuestaPersonaPadron, error) {
	if tipoDocumento != padron.TipoDocumentoDNI &&
		tipoDocumento != padron.TipoDocumentoCUIT &&
		tipoDocumento != padron.TipoDocumentoCUIL {
		return nil, errors.New("el tipo de documento P13 debe ser dni, cuit o cuil")
	}
	ruta := "/api/v2/arca/padron/p13/personas/" +
		url.PathEscape(string(tipoDocumento)) + "/" + url.PathEscape(identificador)
	return s.consultar(ctx, ruta)
}

func (s *ServicioPadron) consultar(ctx context.Context, ruta string) (*padron.RespuestaPersonaPadron, error) {
	var respuesta padron.RespuestaPersonaPadron
	if err := s.transporte.json(ctx, http.MethodGet, ruta, nil, nil, true, &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}
