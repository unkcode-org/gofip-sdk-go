// Modelos del contrato OpenAPI 2.18.0, schemas RespuestaFacturaPersistida y DatosFacturaCompleta.
package facturacion

import "time"

type RespuestaFacturaPersistida struct {
	DocumentoID       string               `json:"documento_id"`
	EmisorID          string               `json:"emisor_id"`
	PuntoVentaID      string               `json:"punto_venta_id"`
	Estado            string               `json:"estado"`
	Origen            string               `json:"origen"`
	ExternalReference string               `json:"external_reference,omitempty"`
	Factura           *DetalleDatosFactura `json:"factura"`
	PDF               EstadoPDFFactura     `json:"pdf"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
}

type EstadoPDFFactura struct {
	Disponible    bool       `json:"disponible"`
	ContentType   string     `json:"content_type,omitempty"`
	GeneradoAt    *time.Time `json:"generado_at,omitempty"`
	ActualizadoAt *time.Time `json:"actualizado_at,omitempty"`
}

type DetalleDatosFactura struct {
	FacturaID             string                       `json:"factura_id"`
	CAE                   string                       `json:"cae"`
	FechaVencimientoCAE   string                       `json:"fecha_vencimiento_cae"`
	FechaComprobante      string                       `json:"fecha_comprobante"`
	TipoComprobante       int                          `json:"tipo_comprobante"`
	NumeroComprobante     int64                        `json:"numero_comprobante"`
	Concepto              int                          `json:"concepto"`
	PuntoVenta            int                          `json:"punto_venta"`
	TipoDocumento         int                          `json:"tipo_documento"`
	NumeroDocumento       string                       `json:"numero_documento"`
	MonedaID              string                       `json:"moneda_id"`
	Cotizacion            float64                      `json:"cotizacion"`
	ImporteTotal          float64                      `json:"importe_total"`
	ImporteNoGravado      float64                      `json:"importe_no_gravado"`
	ImporteNeto           float64                      `json:"importe_neto"`
	ImporteExento         float64                      `json:"importe_exento"`
	ImporteTributos       float64                      `json:"importe_tributos"`
	ImporteIVA            float64                      `json:"importe_iva"`
	Items                 []DetalleItemFactura         `json:"items"`
	ServicioDetalle       *DetalleServicioFactura      `json:"servicio_detalle"`
	ComprobantesAsociados []DetalleComprobanteAsociado `json:"comprobantes_asociados"`
	AlicuotasIVA          []DetalleIVAFactura          `json:"alicuotas_iva"`
	Tributos              []DetalleTributoFactura      `json:"tributos"`
	CreatedAt             time.Time                    `json:"created_at"`
}

type DetalleItemFactura struct {
	ID           string  `json:"id"`
	Descripcion  string  `json:"descripcion"`
	CodigoIVA    string  `json:"codigo_iva"`
	IVAExento    bool    `json:"iva_exento"`
	IVANoGravado bool    `json:"iva_no_gravado"`
	Precio       float64 `json:"precio"`
	Cantidad     float64 `json:"cantidad"`
	Subtotal     float64 `json:"subtotal"`
	MontoIVA     float64 `json:"monto_iva"`
	Total        float64 `json:"total"`
}

type DetalleServicioFactura struct {
	FechaVencimientoPago string `json:"fecha_vencimiento_pago"`
	FechaServicioDesde   string `json:"fecha_servicio_desde"`
	FechaServicioHasta   string `json:"fecha_servicio_hasta"`
}

type DetalleComprobanteAsociado struct {
	Tipo       int    `json:"tipo"`
	PuntoVenta int    `json:"punto_venta"`
	Numero     int    `json:"numero"`
	CUIT       string `json:"cuit"`
	Fecha      string `json:"fecha"`
}

type DetalleIVAFactura struct {
	CodigoIVA int     `json:"codigo_iva"`
	BaseImp   float64 `json:"base_imponible"`
	Importe   float64 `json:"importe"`
}

type DetalleTributoFactura struct {
	ID          int     `json:"id"`
	Descripcion string  `json:"descripcion"`
	BaseImp     float64 `json:"base_imponible"`
	Alicuota    float64 `json:"alicuota"`
	Importe     float64 `json:"importe"`
}
