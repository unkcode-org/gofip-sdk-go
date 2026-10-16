package consumer_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	gofip "github.com/unkcode-org/gofip-sdk-go/v2"
	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/documentos"
	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
	sepmodel "github.com/unkcode-org/gofip-sdk-go/v2/modelos/sep"
)

// Consume la API pública desde un módulo separado y verifica el contrato HTTP
// sin credenciales reales ni efectos fiscales.
func TestConsumidorV2(t *testing.T) {
	const documentoID = "01900000-0000-7000-8000-000000000001"
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if strings.HasPrefix(r.URL.Path, "/api/v2/embedded/") {
			if r.Header.Get("Authorization") != "IntegrationKey sep-prueba" || r.Header.Get("X-ARCA-Ambiente") != "" {
				t.Error("SEP incorrecto")
			}
			w.WriteHeader(201)
			if strings.HasSuffix(r.URL.Path, "/sessions") {
				_, _ = fmt.Fprintf(w, `{"data":{"document_url":"https://embedded.example/capacidad","document_expires_at":"2026-10-16T13:00:00Z","bootstrap_token":"gofip_emb_%s","bootstrap_expires_at":"2026-10-16T12:01:00Z"}}`, strings.Repeat("a", 43))
			} else {
				_, _ = fmt.Fprintf(w, `{"data":{"bootstrap_token":"gofip_emb_%s","bootstrap_expires_at":"2026-10-16T12:01:00Z"}}`, strings.Repeat("a", 43))
			}
			return
		}
		if r.Header.Get("Authorization") != "Bearer api-prueba" {
			t.Error("Bearer incorrecto")
		}
		switch r.URL.Path {
		case "/api/v2/arca/factura":
			if r.Header.Get("Idempotency-Key") != "venta-1" || r.Header.Get("X-External-Reference") != "venta-1" {
				t.Error("opciones existentes incorrectas")
			}
			_, _ = fmt.Fprintf(w, `{"documento_id":%q,"factura":{"cae":"00123456789012"}}`, documentoID)
		case "/api/v2/arca/comprobantes/1/11/1":
			_, _ = fmt.Fprint(w, `{"codigo_autorizacion":"00123456789012"}`)
		case "/api/v2/templates-oficiales/recibo/documentos":
			var cuerpo map[string]json.RawMessage
			if json.NewDecoder(r.Body).Decode(&cuerpo) != nil || cuerpo["data"] == nil {
				t.Error("recibo sin data")
			}
			_, _ = fmt.Fprint(w, `{"documento_id":"recibo-1","pdf_generated_at":"2026-10-16T12:00:00Z"}`)
		default:
			_, _ = fmt.Fprintf(w, `{"documento_id":%q,"estado":"autorizado","factura":{"cae":"00123456789012"},"pdf":{"disponible":false},"created_at":"2026-10-16T12:00:00Z","updated_at":"2026-10-16T12:00:00Z"}`, documentoID)
		}
	}))
	defer s.Close()
	c, err := gofip.NuevoCliente(gofip.OpcionesCliente{ClaveAPI: "api-prueba", URLBase: s.URL, AmbienteARCA: gofip.AmbienteHomologacion, ClienteHTTP: s.Client()})
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if _, err := c.Facturacion.EmitirFactura(ctx, facturacion.SolicitudEmitirFactura{}, gofip.OpcionesEmisionFactura{ClaveIdempotencia: "venta-1", ReferenciaExterna: "venta-1"}); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Facturacion.ObtenerComprobante(ctx, 1, facturacion.ComprobanteFacturaC, 1); err != nil {
		t.Fatal(err)
	}
	porCAE, err := c.Facturacion.ObtenerFacturaPorCAE(ctx, "00123456789012")
	if err != nil {
		t.Fatal(err)
	}
	porID, err := c.Facturacion.ObtenerFacturaPorId(ctx, documentoID)
	if err != nil || porID.DocumentoID != porCAE.DocumentoID {
		t.Fatal("consultas incompatibles")
	}
	if _, err := c.Documentos.GenerarRecibo(ctx, documentos.DatosRecibo{}, gofip.OpcionesGenerarRecibo{}); err != nil {
		t.Fatal(err)
	}
	sep, err := gofip.NuevoClienteSEP(gofip.OpcionesClienteSEP{ClaveIntegracion: "sep-prueba", URLBase: s.URL, ClienteHTTP: s.Client()})
	if err != nil {
		t.Fatal(err)
	}
	solicitud := sepmodel.SolicitudCrearSesionEmbedded{IDEmisorExterno: "elixir", Permisos: []sepmodel.PermisoEmbedded{sepmodel.PermisoFacturasLeer}}
	sesion, err := sep.CrearSesion(ctx, solicitud)
	if err != nil {
		t.Fatal(err)
	}
	if sesion.Datos.TokenBootstrap.Valor() == "" {
		t.Fatal("capacidad ausente")
	}
	if _, err := sep.CrearBootstrap(ctx, solicitud); err != nil {
		t.Fatal(err)
	}
}
