package gofip

import (
	"context"
	"fmt"
	"net/http"

	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
)

type ServicioFacturacion struct {
	transporte *transporte
}

type OpcionesEmisionFactura struct {
	ClaveIdempotencia string
	ReferenciaExterna string
}

func (s *ServicioFacturacion) EmitirFactura(
	ctx context.Context,
	solicitud facturacion.SolicitudEmitirFactura,
	opciones OpcionesEmisionFactura,
) (*facturacion.RespuestaEmitirFactura, error) {
	var respuesta facturacion.RespuestaEmitirFactura
	cabeceras := map[string]string{
		"Idempotency-Key":      opciones.ClaveIdempotencia,
		"X-External-Reference": opciones.ReferenciaExterna,
	}
	if err := s.transporte.json(
		ctx,
		http.MethodPost,
		"/api/v2/arca/factura",
		cuerpoDatos[facturacion.SolicitudEmitirFactura]{Datos: solicitud},
		cabeceras,
		true,
		&respuesta,
	); err != nil {
		return nil, err
	}
	return &respuesta, nil
}

func (s *ServicioFacturacion) ObtenerComprobante(
	ctx context.Context,
	puntoVenta int,
	tipoComprobante facturacion.TipoComprobante,
	numeroComprobante int64,
) (*facturacion.ComprobanteConsultado, error) {
	ruta := fmt.Sprintf(
		"/api/v2/arca/comprobantes/%d/%d/%d",
		puntoVenta,
		tipoComprobante,
		numeroComprobante,
	)
	var respuesta facturacion.ComprobanteConsultado
	if err := s.transporte.json(ctx, http.MethodGet, ruta, nil, nil, true, &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}
