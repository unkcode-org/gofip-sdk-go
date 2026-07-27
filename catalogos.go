package gofip

import (
	"context"
	"net/http"
	"net/url"

	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/catalogos"
)

type ServicioCatalogos struct {
	transporte *transporte
}

func (s *ServicioCatalogos) Estado(ctx context.Context) (*catalogos.RespuestaEstadoServicio, error) {
	var respuesta catalogos.RespuestaEstadoServicio
	if err := s.obtener(ctx, "/api/v2/arca/estado", &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}

func (s *ServicioCatalogos) ListarPuntosVenta(ctx context.Context) (*catalogos.RespuestaPuntosVentaARCA, error) {
	var respuesta catalogos.RespuestaPuntosVentaARCA
	if err := s.obtener(ctx, "/api/v2/arca/puntos-venta", &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}

func (s *ServicioCatalogos) ListarTributos(ctx context.Context) (*catalogos.RespuestaTributos, error) {
	var respuesta catalogos.RespuestaTributos
	if err := s.obtener(ctx, "/api/v2/arca/tributos", &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}

func (s *ServicioCatalogos) ListarConceptos(ctx context.Context) (*catalogos.RespuestaConceptos, error) {
	var respuesta catalogos.RespuestaConceptos
	if err := s.obtener(ctx, "/api/v2/arca/conceptos", &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}

func (s *ServicioCatalogos) ListarAlicuotas(ctx context.Context) (*catalogos.RespuestaAlicuotas, error) {
	var respuesta catalogos.RespuestaAlicuotas
	if err := s.obtener(ctx, "/api/v2/arca/alicuotas", &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}

func (s *ServicioCatalogos) ListarTiposComprobante(ctx context.Context) (*catalogos.RespuestaTiposComprobante, error) {
	var respuesta catalogos.RespuestaTiposComprobante
	if err := s.obtener(ctx, "/api/v2/arca/tipos-comprobante", &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}

func (s *ServicioCatalogos) ListarOpcionales(ctx context.Context) (*catalogos.RespuestaOpcionales, error) {
	var respuesta catalogos.RespuestaOpcionales
	if err := s.obtener(ctx, "/api/v2/arca/opcionales", &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}

func (s *ServicioCatalogos) ListarTiposDocumento(ctx context.Context) (*catalogos.RespuestaTiposDocumento, error) {
	var respuesta catalogos.RespuestaTiposDocumento
	if err := s.obtener(ctx, "/api/v2/arca/tipos-documento", &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}

func (s *ServicioCatalogos) ListarMonedas(ctx context.Context) (*catalogos.RespuestaMonedas, error) {
	var respuesta catalogos.RespuestaMonedas
	if err := s.obtener(ctx, "/api/v2/arca/monedas", &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}

func (s *ServicioCatalogos) ObtenerCotizacion(ctx context.Context, monedaID string) (*catalogos.RespuestaCotizacionMoneda, error) {
	var respuesta catalogos.RespuestaCotizacionMoneda
	ruta := "/api/v2/arca/monedas/" + url.PathEscape(monedaID) + "/cotizacion"
	if err := s.obtener(ctx, ruta, &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}

func (s *ServicioCatalogos) obtener(ctx context.Context, ruta string, destino any) error {
	return s.transporte.json(ctx, http.MethodGet, ruta, nil, nil, true, destino)
}
