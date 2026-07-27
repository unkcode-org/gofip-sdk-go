package facturacion

import "github.com/unkcode-org/gofip-sdk-go/v2/modelos/catalogos"

type TipoComprobante int
type Concepto int
type TipoDocumento int
type CondicionFiscal int
type AlicuotaIVA string
type TamanoPapel string

const (
	TamanoPapelA4   TamanoPapel = "a4"
	TamanoPapel80MM TamanoPapel = "80mm"

	CondicionResponsableMonotributo    CondicionFiscal = 1
	CondicionSujetoNoCategorizado      CondicionFiscal = 2
	CondicionConsumidorFinal           CondicionFiscal = 3
	CondicionProveedorExterior         CondicionFiscal = 4
	CondicionIVAExento                 CondicionFiscal = 5
	CondicionIVAResponsableInscripto   CondicionFiscal = 6
	CondicionIVAResponsableNoInscripto CondicionFiscal = 7

	ComprobanteFacturaA TipoComprobante = 1
	ComprobanteFacturaB TipoComprobante = 6
	ComprobanteFacturaC TipoComprobante = 11

	ConceptoProductos Concepto = 1
	ConceptoServicios Concepto = 2
	ConceptoAmbos     Concepto = 3

	DocumentoCUIT            TipoDocumento = 80
	DocumentoCUIL            TipoDocumento = 86
	DocumentoDNI             TipoDocumento = 96
	DocumentoConsumidorFinal TipoDocumento = 99

	AlicuotaCero         AlicuotaIVA = "3"
	AlicuotaDosConCinco  AlicuotaIVA = "9"
	AlicuotaCinco        AlicuotaIVA = "8"
	AlicuotaDiezConCinco AlicuotaIVA = "4"
	AlicuotaVeintiuno    AlicuotaIVA = "5"
	AlicuotaVeintisiete  AlicuotaIVA = "6"
)

type Comprobante struct {
	TipoComprobante       TipoComprobante        `json:"tipo_comprobante"`
	Concepto              Concepto               `json:"concepto"`
	TipoDocumento         TipoDocumento          `json:"tipo_documento"`
	NumeroDocumento       int64                  `json:"numero_documento"`
	CondicionFiscal       CondicionFiscal        `json:"condicion_fiscal"`
	FechaComprobante      *string                `json:"fecha_comprobante,omitempty"`
	MonedaID              *string                `json:"moneda_id,omitempty"`
	Cotizacion            *float64               `json:"cotizacion,omitempty"`
	FechaServicioDesde    *string                `json:"fecha_servicio_desde,omitempty"`
	FechaServicioHasta    *string                `json:"fecha_servicio_hasta,omitempty"`
	FechaVencimientoPago  *string                `json:"fecha_vencimiento_pago,omitempty"`
	CancelaMismaMoneda    *string                `json:"cancela_misma_moneda,omitempty"`
	Detalles              []ItemComprobante      `json:"detalles"`
	ComprobantesAsociados []ComprobanteAsociado  `json:"comprobantes_asociados,omitempty"`
	Tributos              []TributoComprobante   `json:"tributos,omitempty"`
	Opcionales            []OpcionalComprobante  `json:"opcionales,omitempty"`
	Compradores           []CompradorComprobante `json:"compradores,omitempty"`
	ConsolidarItems       *bool                  `json:"consolidar_items,omitempty"`
}

type ItemComprobante struct {
	ID             string      `json:"id"`
	Descripcion    string      `json:"descripcion"`
	PrecioUnitario float64     `json:"precio_unitario"`
	Cantidad       float64     `json:"cantidad"`
	AlicuotasIVA   AlicuotaIVA `json:"alicuotas_iva"`
	IVAExento      *bool       `json:"iva_exento,omitempty"`
	NoGravadoIVA   *bool       `json:"no_gravadode_iva,omitempty"`
	Subtotal       float64     `json:"subtotal"`
	MontoIVA       *float64    `json:"monto_iva,omitempty"`
	TotalConIVA    *float64    `json:"total_con_iva,omitempty"`
	PrecioConIVA   *float64    `json:"precio_con_iva,omitempty"`
}

type IVAComprobante struct {
	ID            int     `json:"id"`
	BaseImponible float64 `json:"base_imponible"`
	Importe       float64 `json:"importe"`
}

type TributoComprobante struct {
	ID            int     `json:"id"`
	Descripcion   string  `json:"descripcion"`
	BaseImponible float64 `json:"base_imponible"`
	Alicuota      float64 `json:"alicuota"`
	Importe       float64 `json:"importe"`
}

type OpcionalComprobante struct {
	ID    int     `json:"id"`
	Valor float64 `json:"valor"`
}

type CompradorComprobante struct {
	TipoDocumento   int     `json:"tipo_documento"`
	NumeroDocumento int64   `json:"numero_documento"`
	Porcentaje      float64 `json:"porcentaje"`
}

type ComprobanteAsociado struct {
	Numero     int   `json:"numero"`
	Tipo       int   `json:"tipo"`
	CUIT       int64 `json:"cuit"`
	PuntoVenta int   `json:"punto_venta"`
}

type ActividadComprobante struct {
	ID int64 `json:"id"`
}

type PeriodoComprobante struct {
	FechaDesde string `json:"fecha_desde"`
	FechaHasta string `json:"fecha_hasta"`
}

type ObservacionComprobante struct {
	Codigo  *int    `json:"codigo,omitempty"`
	Mensaje *string `json:"mensaje,omitempty"`
}

type MetadatosComprobante struct {
	CondicionFiscalCliente    CondicionFiscal       `json:"condicion_fiscal_cliente"`
	MetodoPago                catalogos.MedioPagoID `json:"metodo_pago"`
	TelefonoCliente           *string               `json:"telefono_cliente,omitempty"`
	EmailCliente              *string               `json:"email_cliente,omitempty"`
	DomicilioCliente          *string               `json:"domicilio_cliente,omitempty"`
	NombreORazonSocialCliente *string               `json:"nombre_o_razon_social_cliente,omitempty"`
	TotalOINI                 *float64              `json:"total_oini,omitempty"`
	AliasMoneda               *string               `json:"alias_moneda,omitempty"`
}

type SolicitudEmitirFactura struct {
	Factura     Comprobante           `json:"factura"`
	Metadatos   *MetadatosComprobante `json:"metadatos,omitempty"`
	TamanoPapel *TamanoPapel          `json:"tamano_papel,omitempty"`
}

type ComprobanteAutorizado struct {
	CAE                   *string                `json:"cae,omitempty"`
	FechaVencimientoCAE   *string                `json:"fecha_vencimiento_cae,omitempty"`
	FechaComprobante      *string                `json:"fecha_comprobante,omitempty"`
	TipoComprobante       *int                   `json:"tipo_comprobante,omitempty"`
	PuntoVenta            *int                   `json:"punto_venta,omitempty"`
	NumeroComprobante     *int64                 `json:"numero_comprobante,omitempty"`
	Concepto              *int                   `json:"concepto,omitempty"`
	TipoDocumento         *int                   `json:"tipo_documento,omitempty"`
	NumeroDocumento       *int64                 `json:"numero_documento,omitempty"`
	MonedaID              *string                `json:"moneda_id,omitempty"`
	Cotizacion            *float64               `json:"cotizacion,omitempty"`
	ImporteTotal          *float64               `json:"importe_total,omitempty"`
	ImporteNoGravado      *float64               `json:"importe_no_gravado,omitempty"`
	ImporteNeto           *float64               `json:"importe_neto,omitempty"`
	ImporteExento         *float64               `json:"importe_exento,omitempty"`
	ImporteTributos       *float64               `json:"importe_tributos,omitempty"`
	ImporteIVA            *float64               `json:"importe_iva,omitempty"`
	FechaVencimientoPago  *string                `json:"fecha_vencimiento_pago,omitempty"`
	FechaServicioDesde    *string                `json:"fecha_servicio_desde,omitempty"`
	FechaServicioHasta    *string                `json:"fecha_servicio_hasta,omitempty"`
	AlicuotasIVA          []IVAComprobante       `json:"alicuotas_iva,omitempty"`
	Periodo               *PeriodoComprobante    `json:"periodo,omitempty"`
	Tributos              []TributoComprobante   `json:"tributos,omitempty"`
	ComprobantesAsociados []ComprobanteAsociado  `json:"comprobantes_asociados,omitempty"`
	Opcionales            []OpcionalComprobante  `json:"opcionales,omitempty"`
	Compradores           []CompradorComprobante `json:"compradores,omitempty"`
	Actividades           []ActividadComprobante `json:"actividades,omitempty"`
}

type RespuestaEmitirFactura struct {
	DocumentoID    *string               `json:"documento_id,omitempty"`
	Factura        ComprobanteAutorizado `json:"factura"`
	PDFURL         *string               `json:"pdf_url,omitempty"`
	PDFContentType *string               `json:"pdf_content_type,omitempty"`
}

type ComprobanteConsultado struct {
	PuntoVenta             *int                     `json:"punto_venta,omitempty"`
	Resultado              *string                  `json:"resultado,omitempty"`
	TipoEmision            *string                  `json:"tipo_emision,omitempty"`
	CodigoAutorizacion     *string                  `json:"codigo_autorizacion,omitempty"`
	TipoDocumento          *int                     `json:"tipo_documento,omitempty"`
	NumeroDocumento        *int64                   `json:"numero_documento,omitempty"`
	Concepto               *int                     `json:"concepto,omitempty"`
	TipoComprobante        *int                     `json:"tipo_comprobante,omitempty"`
	ComprobanteDesde       *int64                   `json:"comprobante_desde,omitempty"`
	ComprobanteHasta       *int64                   `json:"comprobante_hasta,omitempty"`
	FechaComprobante       *string                  `json:"fecha_comprobante,omitempty"`
	ImporteTotal           *float64                 `json:"importe_total,omitempty"`
	ImporteNoGravado       *float64                 `json:"importe_no_gravado,omitempty"`
	ImporteNeto            *float64                 `json:"importe_neto,omitempty"`
	ImporteExento          *float64                 `json:"importe_exento,omitempty"`
	ImporteTributos        *float64                 `json:"importe_tributos,omitempty"`
	ImporteIVA             *float64                 `json:"importe_iva,omitempty"`
	FechaVencimiento       *string                  `json:"fecha_vencimiento,omitempty"`
	FechaProceso           *string                  `json:"fecha_proceso,omitempty"`
	FechaVencimientoPago   *string                  `json:"fecha_vencimiento_pago,omitempty"`
	FechaServicioDesde     *string                  `json:"fecha_servicio_desde,omitempty"`
	FechaServicioHasta     *string                  `json:"fecha_servicio_hasta,omitempty"`
	MonedaID               *string                  `json:"moneda_id,omitempty"`
	Cotizacion             *float64                 `json:"cotizacion,omitempty"`
	CancelaMismaMoneda     *string                  `json:"cancela_misma_moneda,omitempty"`
	IDCondicionIVAReceptor *int                     `json:"id_condicion_iva_receptor,omitempty"`
	AlicuotasIVA           []IVAComprobante         `json:"alicuotas_iva,omitempty"`
	Tributos               []TributoComprobante     `json:"tributos,omitempty"`
	ComprobantesAsociados  []ComprobanteAsociado    `json:"comprobantes_asociados,omitempty"`
	Opcionales             []OpcionalComprobante    `json:"opcionales,omitempty"`
	Compradores            []CompradorComprobante   `json:"compradores,omitempty"`
	Actividades            []ActividadComprobante   `json:"actividades,omitempty"`
	PeriodoAsociado        *PeriodoComprobante      `json:"periodo_asociado,omitempty"`
	Observaciones          []ObservacionComprobante `json:"observaciones,omitempty"`
}
