package consumer_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	gofip "github.com/unkcode-org/gofip-sdk-go/v2"
	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
)

func TestConsumidorImportacionHistorica(t *testing.T) {
	const id = "01900000-0000-7000-8000-000000000001"
	calls := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Path != "/api/v2/facturas/importaciones" || r.Method != "POST" || r.Header.Get("X-ARCA-Ambiente") != "" || r.Header.Get("Idempotency-Key") != "original" {
			t.Error("ruta mutante o contrato incorrecto")
		}
		w.WriteHeader(200)
		_, _ = fmt.Fprintf(w, `{"documento_id":%q,"importado":false,"ya_existia":true,"verificacion":{"metodo":"wsfe.FECompConsultar","verificado_at":"2026-10-16T20:00:00Z"},"documento":{"documento_id":%q,"emisor_id":"emisor-a","estado":"autorizado","origen":"api","factura":{"cae":"00123456789012","punto_venta":3,"tipo_comprobante":6,"numero_comprobante":12345}}}`, id, id)
	}))
	defer s.Close()
	c, err := gofip.NuevoCliente(gofip.OpcionesCliente{URLBase: s.URL, ClaveAPI: "test", AmbienteARCA: gofip.AmbienteHomologacion, ClienteHTTP: s.Client()})
	if err != nil {
		t.Fatal(err)
	}
	r, err := c.Facturacion.ImportarFactura(context.Background(), facturacion.SolicitudImportarFactura{Factura: facturacion.FacturaHistorica{PuntoVenta: 3, TipoComprobante: facturacion.TipoComprobante(6), NumeroComprobante: 12345, CAE: "00123456789012"}, Procedencia: facturacion.ProcedenciaFactura{Sistema: "gofip_legacy"}}, gofip.OpcionesImportacionFactura{ClaveIdempotencia: "original", ReferenciaExterna: "cliente:original"})
	if err != nil || r.Importado || !r.YaExistia || r.DocumentoID != id || r.Documento.Origen != "api" || calls != 1 {
		t.Fatalf("importación incompatible: %+v %v", r, err)
	}
}
