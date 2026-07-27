package documentos

import (
	"time"

	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
)

const (
	SlugRecibo    = "recibo"
	VersionRecibo = "1.0.0"
)

type SolicitudDocumentoComprobante struct {
	Comprobante facturacion.ComprobanteAutorizado `json:"comprobante"`
	Metadatos   facturacion.MetadatosComprobante  `json:"metadatos"`
	Items       []facturacion.ItemComprobante     `json:"items"`
	TamanoPapel *facturacion.TamanoPapel          `json:"tamano_papel,omitempty"`
}

type SolicitudGenerarPDFFactura struct {
	Metadatos   facturacion.MetadatosComprobante `json:"metadatos"`
	TamanoPapel *facturacion.TamanoPapel         `json:"tamano_papel,omitempty"`
}

type RespuestaPDFFactura struct {
	DocumentoID    string  `json:"documento_id"`
	PDFURL         string  `json:"pdf_url"`
	PDFContentType *string `json:"pdf_content_type,omitempty"`
}

type TemplateOficial struct {
	Slug               string         `json:"slug"`
	Version            string         `json:"version"`
	Nombre             string         `json:"nombre"`
	Descripcion        string         `json:"descripcion"`
	Autor              string         `json:"autor"`
	UltimaModificacion time.Time      `json:"ultima_modificacion"`
	SchemaDatos        map[string]any `json:"schema_datos"`
}

type RespuestaTemplatesOficiales struct {
	Templates []TemplateOficial `json:"templates"`
}

type Recibo struct {
	NumeroRecibo string  `json:"numero_recibo"`
	FechaRecibo  string  `json:"fecha_recibo"`
	Subtotal     float64 `json:"subtotal"`
	Total        float64 `json:"total"`
}

type MetadatosRecibo struct {
	CondicionFiscalCliente    facturacion.CondicionFiscal `json:"condicion_fiscal_cliente"`
	DocumentoCliente          int64                       `json:"documento_cliente"`
	NombreORazonSocialCliente string                      `json:"nombre_o_razon_social_cliente"`
	DomicilioCliente          string                      `json:"domicilio_cliente"`
	TelefonoCliente           *string                     `json:"telefono_cliente,omitempty"`
	EmailCliente              *string                     `json:"email_cliente,omitempty"`
}

type ItemRecibo struct {
	ID          string  `json:"id"`
	Descripcion string  `json:"descripcion"`
	Cantidad    float64 `json:"cantidad"`
	Precio      float64 `json:"precio"`
	Total       float64 `json:"total"`
}

type DatosRecibo struct {
	Recibo    Recibo          `json:"recibo"`
	Metadatos MetadatosRecibo `json:"metadatos"`
	Items     []ItemRecibo    `json:"items"`
}

type SolicitudGenerarRecibo struct {
	Datos DatosRecibo `json:"datos"`
}

type RespuestaDocumentoNoFiscal struct {
	DocumentoID     string    `json:"documento_id"`
	TemplateSlug    string    `json:"template_slug"`
	TemplateVersion string    `json:"template_version"`
	PDFURL          string    `json:"pdf_url"`
	PDFContentType  string    `json:"pdf_content_type"`
	PDFGeneratedAt  time.Time `json:"pdf_generated_at"`
}
