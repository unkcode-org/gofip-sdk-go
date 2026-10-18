package gofip

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
)

// OpcionesImportacionFactura has no environment selector. Legacy imports are
// verified in Production by the server, independently of OpcionesCliente.
type OpcionesImportacionFactura struct {
	ClaveIdempotencia string
	ReferenciaExterna string
}

const limiteImportacionRespuesta = 2 << 20

var patronCodigoImportacion = regexp.MustCompile(`^[0-9]{5}$`)
var patronCAEImportacion = regexp.MustCompile(`^[0-9]{1,20}$`)
var patronDocumentoImportacion = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// ImportarFactura imports an existing authorization; it never calls EmitirFactura.
// The API key may import only invoices from its own issuer and point of sale. No
// application retries, redirects, PDF generation or Legacy fallbacks.
func (s *ServicioFacturacion) ImportarFactura(ctx context.Context, solicitud facturacion.SolicitudImportarFactura, o OpcionesImportacionFactura) (*facturacion.RespuestaImportarFactura, error) {
	f := solicitud.Factura
	if f.PuntoVenta < 1 || f.PuntoVenta > 99999 || f.TipoComprobante < 1 || f.NumeroComprobante < 1 || f.NumeroComprobante > 99999999 || !patronCAEImportacion.MatchString(f.CAE) || solicitud.Procedencia.Sistema != "gofip_legacy" || len(solicitud.Procedencia.RegistroID) > 120 || len(o.ClaveIdempotencia) == 0 || len(o.ClaveIdempotencia) > 120 || len(o.ReferenciaExterna) > 120 {
		return nil, errors.New("solicitud de importación inválida")
	}
	for _, v := range []string{o.ClaveIdempotencia, o.ReferenciaExterna} {
		if strings.IndexFunc(v, func(r rune) bool { return r < 32 || r == 127 }) >= 0 {
			return nil, errors.New("cabecera de importación inválida")
		}
	}
	b, err := json.Marshal(cuerpoDatos[facturacion.SolicitudImportarFactura]{Datos: solicitud})
	if err != nil {
		return nil, errors.New("serializar solicitud de importación inválida")
	}
	if len(b) > 1<<20 {
		return nil, errors.New("solicitud de importación excede 1 MiB")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.transporte.urlBase+"/api/v2/facturas/importaciones", bytes.NewReader(b))
	if err != nil {
		return nil, errors.New("crear solicitud de importación")
	}
	// Disable transparent body replay by the standard transport. A caller may
	// explicitly repeat the operation using its durable fiscal/idempotency keys.
	req.GetBody = nil
	req.Header.Set("Authorization", "Bearer "+s.transporte.claveAPI)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Idempotency-Key", o.ClaveIdempotencia)
	if o.ReferenciaExterna != "" {
		req.Header.Set("X-External-Reference", o.ReferenciaExterna)
	}
	client := *s.transporte.cliente
	client.Jar = nil
	client.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	resp, err := client.Do(req)
	if err != nil {
		return nil, errorImportacionSeguro(ctx, err, "fallo de transporte de importación")
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, limiteImportacionRespuesta+1))
	if err != nil {
		return nil, errorImportacionSeguro(ctx, err, "fallo al leer respuesta de importación")
	}
	if len(raw) > limiteImportacionRespuesta {
		return nil, errors.New("respuesta de importación excede 2 MiB")
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var problem ErrorAPI
		if json.Unmarshal(raw, &problem) == nil && patronCodigoImportacion.MatchString(problem.Codigo) {
			return nil, &ErrorAPI{Estado: resp.StatusCode, Codigo: problem.Codigo, Titulo: "Error de importación fiscal", Detalle: http.StatusText(resp.StatusCode)}
		}
		return nil, &ErrorHTTP{Estado: resp.StatusCode, Cuerpo: http.StatusText(resp.StatusCode)}
	}
	var result facturacion.RespuestaImportarFactura
	dec := json.NewDecoder(bytes.NewReader(raw))
	if dec.Decode(&result) != nil {
		return nil, errors.New("respuesta de importación inválida")
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return nil, errors.New("respuesta de importación contiene datos adicionales")
	}
	if !patronDocumentoImportacion.MatchString(result.DocumentoID) || result.Importado == result.YaExistia || result.Importado != (resp.StatusCode == 201) || result.Verificacion.Metodo != "wsfe.FECompConsultar" || result.Verificacion.VerificadoAt.IsZero() || result.Documento == nil || result.Documento.DocumentoID != result.DocumentoID || result.Documento.EmisorID == "" || result.Documento.Estado != "autorizado" || result.Documento.Factura == nil || result.Documento.Factura.CAE != f.CAE || result.Documento.Factura.PuntoVenta != f.PuntoVenta || result.Documento.Factura.TipoComprobante != int(f.TipoComprobante) || result.Documento.Factura.NumeroComprobante != f.NumeroComprobante {
		return nil, errors.New("respuesta de importación incompleta o incompatible")
	}
	return &result, nil
}
func errorImportacionSeguro(ctx context.Context, err error, message string) error {
	if errors.Is(err, context.Canceled) || errors.Is(ctx.Err(), context.Canceled) {
		return fmt.Errorf("importación fiscal: %w", context.Canceled)
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("importación fiscal: %w", context.DeadlineExceeded)
	}
	// Arbitrary transports may include headers, bodies or URLs in both Error and
	// Unwrap. Preserve only known context causes, never those diagnostic chains.
	return errors.New(message)
}
