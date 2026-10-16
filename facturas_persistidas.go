package gofip

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
)

// OpcionesConsultaFacturaPorCAE reduce la búsqueda dentro del emisor de la API
// key. El punto de venta es su número fiscal, no el UUID interno.
type OpcionesConsultaFacturaPorCAE struct {
	DocumentoID       string
	PuntoVenta        *int
	TipoComprobante   *facturacion.TipoComprobante
	NumeroComprobante *int64
}

// ObtenerFacturaPorCAE lee exclusivamente la base GOFIP. Un CAE ambiguo produce
// ErrorAPI con Estado 409 y Codigo 13006; no consulta ARCA como fallback.
func (s *ServicioFacturacion) ObtenerFacturaPorCAE(ctx context.Context, cae string) (*facturacion.RespuestaFacturaPersistida, error) {
	return s.ObtenerFacturaPorCAEConOpciones(ctx, cae, OpcionesConsultaFacturaPorCAE{})
}

func (s *ServicioFacturacion) ObtenerFacturaPorCAEConOpciones(ctx context.Context, cae string, o OpcionesConsultaFacturaPorCAE) (*facturacion.RespuestaFacturaPersistida, error) {
	if cae == "" {
		return nil, errors.New("el CAE es obligatorio")
	}
	q := url.Values{}
	if o.DocumentoID != "" {
		q.Set("documento_id", o.DocumentoID)
	}
	if o.PuntoVenta != nil {
		q.Set("punto_venta", strconv.Itoa(*o.PuntoVenta))
	}
	if o.TipoComprobante != nil {
		q.Set("tipo_comprobante", strconv.Itoa(int(*o.TipoComprobante)))
	}
	if o.NumeroComprobante != nil {
		q.Set("numero_comprobante", strconv.FormatInt(*o.NumeroComprobante, 10))
	}
	ruta := "/api/v2/facturas/cae/" + url.PathEscape(cae)
	if len(q) > 0 {
		ruta += "?" + q.Encode()
	}
	return s.obtenerFacturaPersistida(ctx, ruta)
}

// ObtenerFacturaPorId usa el ID estable del documento GOFIP. Permite leer
// documentos de todos los puntos de venta del mismo emisor, nunca de otro.
// No selecciona ambiente: ese dato no está persistido en las facturas actuales.
func (s *ServicioFacturacion) ObtenerFacturaPorId(ctx context.Context, documentoID string) (*facturacion.RespuestaFacturaPersistida, error) {
	if documentoID == "" {
		return nil, errors.New("el documento ID es obligatorio")
	}
	return s.obtenerFacturaPersistida(ctx, "/api/v2/facturas/"+url.PathEscape(documentoID))
}

func (s *ServicioFacturacion) obtenerFacturaPersistida(ctx context.Context, ruta string) (*facturacion.RespuestaFacturaPersistida, error) {
	var respuesta facturacion.RespuestaFacturaPersistida
	if err := s.transporte.json(ctx, http.MethodGet, ruta, nil, nil, false, &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}
