package gofip

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const documentoFiscalPrueba = "01900000-0000-7000-8000-000000000001"

func TestConsultasPersistidasBearerMismoModelo(t *testing.T) {
	var rutas []string
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rutas = append(rutas, r.URL.Path)
		if r.Method != "GET" || r.Header.Get("Authorization") != "Bearer clave-prueba" {
			t.Error("autenticación o método incorrectos")
		}
		if r.Header.Get("X-ARCA-Ambiente") != "" {
			t.Error("consulta persistida no selecciona ambiente")
		}
		if r.URL.Path == "/api/v2/facturas/cae/00123456789012" && r.URL.Query().Get("documento_id") != documentoFiscalPrueba {
			t.Error("falta filtro")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"documento_id":"`+documentoFiscalPrueba+`","emisor_id":"emisor-a","punto_venta_id":"pv-2","estado":"autorizado","origen":"api","external_reference":"venta-123","factura":{"factura_id":"f-1","cae":"00123456789012","numero_documento":"00123456789","numero_comprobante":125,"items":[{"id":"item-1","descripcion":"Anillo","total":121}],"alicuotas_iva":[{"codigo_iva":5,"base_imponible":100,"importe":21}],"tributos":[],"comprobantes_asociados":[],"created_at":"2026-10-16T12:00:00Z"},"pdf":{"disponible":true},"created_at":"2026-10-16T12:00:00Z","updated_at":"2026-10-16T12:00:00Z"}`)
	}))
	defer servidor.Close()
	cliente := nuevoClientePrueba(t, servidor.URL, servidor.Client())
	factura, err := cliente.Facturacion.ObtenerFacturaPorId(context.Background(), documentoFiscalPrueba)
	if err != nil {
		t.Fatal(err)
	}
	if factura.Factura.CAE != "00123456789012" || factura.Factura.NumeroDocumento != "00123456789" || factura.Factura.Items[0].Total != 121 || factura.CreatedAt.Year() != 2026 {
		t.Fatal("detalle fiscal incorrecto")
	}
	porCAE, err := cliente.Facturacion.ObtenerFacturaPorCAEConOpciones(context.Background(), "00123456789012", OpcionesConsultaFacturaPorCAE{DocumentoID: documentoFiscalPrueba})
	if err != nil || porCAE.DocumentoID != factura.DocumentoID {
		t.Fatalf("consulta CAE: %v", err)
	}
	if len(rutas) != 2 || rutas[0] != "/api/v2/facturas/"+documentoFiscalPrueba {
		t.Fatalf("rutas: %v", rutas)
	}
}

func TestConsultaCAEPropagaConflictoYCancelacion(t *testing.T) {
	llamadas := 0
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		llamadas++
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(409)
		_, _ = io.WriteString(w, `{"code":"13006","title":"CAE ambiguo","status":409,"detail":"Agregue filtros"}`)
	}))
	defer s.Close()
	c := nuevoClientePrueba(t, s.URL, s.Client())
	_, err := c.Facturacion.ObtenerFacturaPorCAE(context.Background(), "00123456789012")
	var api *ErrorAPI
	if !errors.As(err, &api) || api.Estado != 409 || api.Codigo != "13006" {
		t.Fatalf("error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = c.Facturacion.ObtenerFacturaPorId(ctx, documentoFiscalPrueba)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelación: %v", err)
	}
	if llamadas != 1 {
		t.Fatalf("reintentos inesperados: %d", llamadas)
	}
}
