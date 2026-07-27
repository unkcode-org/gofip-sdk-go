package catalogos

type MedioPagoID int

const (
	MedioPagoEfectivo        MedioPagoID = 1
	MedioPagoTarjetaDebito   MedioPagoID = 2
	MedioPagoTarjetaCredito  MedioPagoID = 3
	MedioPagoCheque          MedioPagoID = 4
	MedioPagoCuentaCorriente MedioPagoID = 5
	MedioPagoTicket          MedioPagoID = 6
	MedioPagoTransferencia   MedioPagoID = 7
	MedioPagoOtro            MedioPagoID = 8
)

type MedioPago struct {
	ID          MedioPagoID `json:"id"`
	Descripcion string      `json:"descripcion"`
}

type RespuestaMediosPago struct {
	MediosPago []MedioPago `json:"medios_pago"`
}

type EstadoWS struct {
	ServidorAplicacion    *string `json:"servidor_aplicacion,omitempty"`
	ServidorBaseDatos     *string `json:"servidor_base_datos,omitempty"`
	ServidorAutenticacion *string `json:"servidor_autenticacion,omitempty"`
}

type RespuestaEstadoServicio struct {
	Estado EstadoWS `json:"estado"`
}

type TipoTributo struct {
	ID          *int    `json:"id,omitempty"`
	Descripcion *string `json:"descripcion,omitempty"`
	FechaDesde  *string `json:"fecha_desde,omitempty"`
	FechaHasta  *string `json:"fecha_hasta,omitempty"`
}

type TipoConcepto struct {
	ID          *int    `json:"id,omitempty"`
	Descripcion *string `json:"descripcion,omitempty"`
	FechaDesde  *string `json:"fecha_desde,omitempty"`
	FechaHasta  *string `json:"fecha_hasta,omitempty"`
}

type TipoAlicuota struct {
	ID          *string `json:"id,omitempty"`
	Descripcion *string `json:"descripcion,omitempty"`
	FechaDesde  *string `json:"fecha_desde,omitempty"`
	FechaHasta  *string `json:"fecha_hasta,omitempty"`
}

type TipoComprobanteARCA struct {
	ID          *int    `json:"id,omitempty"`
	Descripcion *string `json:"descripcion,omitempty"`
	FechaDesde  *string `json:"fecha_desde,omitempty"`
	FechaHasta  *string `json:"fecha_hasta,omitempty"`
}

type TipoDocumentoARCA struct {
	ID          *int    `json:"id,omitempty"`
	Descripcion *string `json:"descripcion,omitempty"`
	FechaDesde  *string `json:"fecha_desde,omitempty"`
	FechaHasta  *string `json:"fecha_hasta,omitempty"`
}

type TipoOpcional struct {
	ID          *string `json:"id,omitempty"`
	Descripcion *string `json:"descripcion,omitempty"`
	FechaDesde  *string `json:"fecha_desde,omitempty"`
	FechaHasta  *string `json:"fecha_hasta,omitempty"`
}

type Moneda struct {
	ID          *string `json:"id,omitempty"`
	Descripcion *string `json:"descripcion,omitempty"`
	FechaDesde  *string `json:"fecha_desde,omitempty"`
	FechaHasta  *string `json:"fecha_hasta,omitempty"`
}

type CotizacionMoneda struct {
	MonedaID   *string  `json:"moneda_id,omitempty"`
	Cotizacion *float64 `json:"cotizacion,omitempty"`
	Fecha      *string  `json:"fch_cotiz,omitempty"`
}

type PuntoVentaARCA struct {
	Numero      *int    `json:"numero,omitempty"`
	TipoEmision *string `json:"tipo_emision,omitempty"`
	Bloqueado   *string `json:"bloqueado,omitempty"`
	FechaBaja   *string `json:"fch_baja,omitempty"`
}

type RespuestaTributos struct {
	Tributos []TipoTributo `json:"tributos"`
}

type RespuestaConceptos struct {
	Conceptos []TipoConcepto `json:"conceptos"`
}

type RespuestaAlicuotas struct {
	Alicuotas []TipoAlicuota `json:"alicuotas"`
}

type RespuestaTiposComprobante struct {
	TiposComprobante []TipoComprobanteARCA `json:"tipos_comprobante"`
}

type RespuestaOpcionales struct {
	Opcionales []TipoOpcional `json:"opcionales"`
}

type RespuestaTiposDocumento struct {
	TiposDocumento []TipoDocumentoARCA `json:"tipos_documento"`
}

type RespuestaMonedas struct {
	Monedas []Moneda `json:"monedas"`
}

type RespuestaCotizacionMoneda struct {
	Cotizacion CotizacionMoneda `json:"cotizacion"`
}

type RespuestaPuntosVentaARCA struct {
	PuntosVenta []PuntoVentaARCA `json:"puntos_venta"`
}
