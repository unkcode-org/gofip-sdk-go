// Modelos de importación del contrato OpenAPI 2.19.0. Legacy siempre consulta Producción.
package facturacion

import (
	"encoding/json"
	"errors"
	"regexp"
	"time"
)

// DecimalHistorico conserva un decimal exacto como string JSON.
type DecimalHistorico string

var patronDecimalHistorico = regexp.MustCompile(`^-?[0-9]{1,16}(\.[0-9]{1,6})?$`)

func (d DecimalHistorico) MarshalJSON() ([]byte, error) {
	if !patronDecimalHistorico.MatchString(string(d)) {
		return nil, errors.New("decimal histórico inválido")
	}
	return json.Marshal(string(d))
}
func (d *DecimalHistorico) UnmarshalJSON(b []byte) error {
	var value string
	if json.Unmarshal(b, &value) != nil || !patronDecimalHistorico.MatchString(value) {
		return errors.New("decimal histórico inválido")
	}
	*d = DecimalHistorico(value)
	return nil
}

type FacturaHistorica struct {
	PuntoVenta        int                      `json:"punto_venta"`
	TipoComprobante   TipoComprobante          `json:"tipo_comprobante"`
	NumeroComprobante int64                    `json:"numero_comprobante"`
	CAE               string                   `json:"cae"`
	DatosHistoricos   *DatosFiscalesHistoricos `json:"datos_historicos,omitempty"`
}

type SolicitudImportarFactura struct {
	Factura     FacturaHistorica         `json:"factura"`
	Procedencia ProcedenciaFactura       `json:"procedencia"`
	Detalle     *DetalleHistoricoFactura `json:"detalle,omitempty"`
}

type ProcedenciaFactura struct {
	Sistema    string `json:"sistema"`
	RegistroID string `json:"registro_id,omitempty"`
}

type DatosFiscalesHistoricos struct {
	Resultado             *string                 `json:"resultado,omitempty"`
	TipoEmision           *string                 `json:"tipo_emision,omitempty"`
	CAE                   *string                 `json:"cae,omitempty"`
	PuntoVenta            *int                    `json:"punto_venta,omitempty"`
	TipoComprobante       *int                    `json:"tipo_comprobante,omitempty"`
	ComprobanteDesde      *int64                  `json:"comprobante_desde,omitempty"`
	ComprobanteHasta      *int64                  `json:"comprobante_hasta,omitempty"`
	FechaVencimientoCAE   *string                 `json:"fecha_vencimiento_cae,omitempty"`
	FechaComprobante      *string                 `json:"fecha_comprobante,omitempty"`
	Concepto              *int                    `json:"concepto,omitempty"`
	TipoDocumento         *int                    `json:"tipo_documento,omitempty"`
	NumeroDocumento       *string                 `json:"numero_documento,omitempty"`
	MonedaID              *string                 `json:"moneda_id,omitempty"`
	Cotizacion            *DecimalHistorico       `json:"cotizacion,omitempty"`
	ImporteTotal          *DecimalHistorico       `json:"importe_total,omitempty"`
	ImporteNeto           *DecimalHistorico       `json:"importe_neto,omitempty"`
	ImporteIVA            *DecimalHistorico       `json:"importe_iva,omitempty"`
	ImporteExento         *DecimalHistorico       `json:"importe_exento,omitempty"`
	ImporteNoGravado      *DecimalHistorico       `json:"importe_no_gravado,omitempty"`
	ImporteTributos       *DecimalHistorico       `json:"importe_tributos,omitempty"`
	FechaServicioDesde    *string                 `json:"fecha_servicio_desde,omitempty"`
	FechaServicioHasta    *string                 `json:"fecha_servicio_hasta,omitempty"`
	FechaVencimientoPago  *string                 `json:"fecha_vencimiento_pago,omitempty"`
	FechaProceso          *string                 `json:"fecha_proceso,omitempty"`
	CancelaMismaMoneda    *string                 `json:"cancela_misma_moneda,omitempty"`
	CondicionIVAReceptor  *int                    `json:"condicion_iva_receptor,omitempty"`
	AlicuotasIVA          *[]IVAHistorico         `json:"alicuotas_iva,omitempty"`
	Tributos              *[]TributoHistorico     `json:"tributos,omitempty"`
	ComprobantesAsociados *[]AsociadoHistorico    `json:"comprobantes_asociados,omitempty"`
	Opcionales            *[]OpcionalHistorico    `json:"opcionales,omitempty"`
	Compradores           *[]CompradorHistorico   `json:"compradores,omitempty"`
	Actividades           *[]ActividadHistorica   `json:"actividades,omitempty"`
	PeriodoAsociado       *PeriodoHistorico       `json:"periodo_asociado,omitempty"`
	Observaciones         *[]ObservacionHistorica `json:"observaciones,omitempty"`
}

type IVAHistorico struct {
	ID            int              `json:"id" xml:"Id"`
	BaseImponible DecimalHistorico `json:"base_imponible" xml:"BaseImp"`
	Importe       DecimalHistorico `json:"importe" xml:"Importe"`
}

type TributoHistorico struct {
	ID            int              `json:"id" xml:"Id"`
	Descripcion   *string          `json:"descripcion,omitempty" xml:"Desc"`
	BaseImponible DecimalHistorico `json:"base_imponible" xml:"BaseImp"`
	Alicuota      DecimalHistorico `json:"alicuota" xml:"Alic"`
	Importe       DecimalHistorico `json:"importe" xml:"Importe"`
}

type OpcionalHistorico struct {
	ID    string  `json:"id" xml:"Id"`
	Valor *string `json:"valor,omitempty" xml:"Valor"`
}

type CompradorHistorico struct {
	TipoDocumento   int              `json:"tipo_documento" xml:"DocTipo"`
	NumeroDocumento string           `json:"numero_documento" xml:"DocNro"`
	Porcentaje      DecimalHistorico `json:"porcentaje" xml:"Porcentaje"`
}

type AsociadoHistorico struct {
	Tipo       int     `json:"tipo" xml:"Tipo"`
	PuntoVenta int     `json:"punto_venta" xml:"PtoVta"`
	Numero     int64   `json:"numero" xml:"Nro"`
	CUIT       *string `json:"cuit,omitempty" xml:"Cuit"`
	Fecha      *string `json:"fecha,omitempty" xml:"CbteFch"`
}

type ActividadHistorica struct {
	ID int64 `json:"id" xml:"Id"`
}

type PeriodoHistorico struct {
	FechaDesde string `json:"fecha_desde" xml:"FchDesde"`
	FechaHasta string `json:"fecha_hasta" xml:"FchHasta"`
}

type ObservacionHistorica struct {
	Codigo  int     `json:"codigo" xml:"Code"`
	Mensaje *string `json:"mensaje,omitempty" xml:"Msg"`
}

type DetalleHistoricoFactura struct {
	Items             []ItemHistorico              `json:"items,omitempty"`
	DatosDocumentales *DatosDocumentalesHistoricos `json:"datos_documentales,omitempty"`
}

type ItemHistorico struct {
	ReferenciaOriginal string           `json:"referencia_original,omitempty"`
	Descripcion        string           `json:"descripcion"`
	Cantidad           DecimalHistorico `json:"cantidad"`
	PrecioUnitario     DecimalHistorico `json:"precio_unitario"`
	Subtotal           DecimalHistorico `json:"subtotal"`
	MontoIVA           DecimalHistorico `json:"monto_iva"`
	Total              DecimalHistorico `json:"total"`
	CodigoIVA          string           `json:"codigo_iva,omitempty"`
	IVAExento          *bool            `json:"iva_exento,omitempty"`
	IVANoGravado       *bool            `json:"iva_no_gravado,omitempty"`
}

type DatosDocumentalesHistoricos struct {
	RazonSocialEmisor          *string `json:"razon_social_emisor,omitempty"`
	DomicilioEmisor            *string `json:"domicilio_emisor,omitempty"`
	CondicionFiscalEmisor      *string `json:"condicion_fiscal_emisor,omitempty"`
	FechaInicioActividadEmisor *string `json:"fecha_inicio_actividad_emisor,omitempty"`
	NombreReceptor             *string `json:"nombre_receptor,omitempty"`
	DomicilioReceptor          *string `json:"domicilio_receptor,omitempty"`
	CondicionFiscalReceptor    *string `json:"condicion_fiscal_receptor,omitempty"`
	MetodoPago                 *string `json:"metodo_pago,omitempty"`
	EmailReceptor              *string `json:"email_receptor,omitempty"`
	TelefonoReceptor           *string `json:"telefono_receptor,omitempty"`
}

type VerificacionImportacion struct {
	Metodo       string    `json:"metodo"`
	VerificadoAt time.Time `json:"verificado_at"`
}

type CompletitudDocumental struct {
	TieneItems             bool     `json:"tiene_items"`
	TieneDatosDocumentales bool     `json:"tiene_datos_documentales"`
	PDFRegenerable         bool     `json:"pdf_regenerable"`
	Motivos                []string `json:"motivos"`
}

type RespuestaImportarFactura struct {
	DocumentoID  string                      `json:"documento_id"`
	Importado    bool                        `json:"importado"`
	YaExistia    bool                        `json:"ya_existia"`
	Verificacion VerificacionImportacion     `json:"verificacion"`
	Documento    *RespuestaFacturaPersistida `json:"documento"`
}
