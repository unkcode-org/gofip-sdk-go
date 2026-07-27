package gofip

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/padron"
)

func TestUtilidadesRespetaAutenticacionYContentType(t *testing.T) {
	var llamadas int
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		llamadas++
		switch r.URL.Path {
		case "/api/v2/salud":
			if got := r.Header.Get("Authorization"); got != "" {
				t.Errorf("salud no debe enviar Authorization, recibió %q", got)
			}
			w.Header().Set("Content-Type", "text/plain")
			_, _ = io.WriteString(w, "OK")
		case "/api/v2/disponibilidad":
			if got := r.Header.Get("Authorization"); got != "Bearer clave-prueba" {
				t.Errorf("Authorization inesperado: %q", got)
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = io.WriteString(w, `"pong"`)
		default:
			http.NotFound(w, r)
		}
	}))
	defer servidor.Close()

	cliente := nuevoClientePrueba(t, servidor.URL, servidor.Client())
	salud, err := cliente.Utilidades.Salud(context.Background())
	if err != nil || salud != "OK" {
		t.Fatalf("Salud() = %q, %v", salud, err)
	}
	disponibilidad, err := cliente.Utilidades.Disponibilidad(context.Background())
	if err != nil || disponibilidad != "pong" {
		t.Fatalf("Disponibilidad() = %q, %v", disponibilidad, err)
	}
	if llamadas != 2 {
		t.Fatalf("llamadas = %d, esperadas 2", llamadas)
	}
}

func TestCatalogosUsaRutasYAmbienteV2(t *testing.T) {
	esperadas := []string{
		"/api/v2/arca/estado",
		"/api/v2/arca/puntos-venta",
		"/api/v2/arca/tributos",
		"/api/v2/arca/conceptos",
		"/api/v2/arca/alicuotas",
		"/api/v2/arca/tipos-comprobante",
		"/api/v2/arca/opcionales",
		"/api/v2/arca/tipos-documento",
		"/api/v2/arca/monedas",
		"/api/v2/arca/monedas/DOL/cotizacion",
	}
	indice := 0
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if indice >= len(esperadas) {
			t.Fatalf("llamada adicional a %s", r.URL.Path)
		}
		if r.URL.Path != esperadas[indice] {
			t.Errorf("ruta %d = %s, esperada %s", indice, r.URL.Path, esperadas[indice])
		}
		indice++
		if got := r.Header.Get("Authorization"); got != "Bearer clave-prueba" {
			t.Errorf("Authorization inesperado: %q", got)
		}
		if got := r.Header.Get("X-ARCA-Ambiente"); got != "homologacion" {
			t.Errorf("ambiente inesperado: %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v2/arca/estado":
			_, _ = io.WriteString(w, `{"estado":{}}`)
		case "/api/v2/arca/puntos-venta":
			_, _ = io.WriteString(w, `{"puntos_venta":[]}`)
		case "/api/v2/arca/tributos":
			_, _ = io.WriteString(w, `{"tributos":[]}`)
		case "/api/v2/arca/conceptos":
			_, _ = io.WriteString(w, `{"conceptos":[]}`)
		case "/api/v2/arca/alicuotas":
			_, _ = io.WriteString(w, `{"alicuotas":[]}`)
		case "/api/v2/arca/tipos-comprobante":
			_, _ = io.WriteString(w, `{"tipos_comprobante":[]}`)
		case "/api/v2/arca/opcionales":
			_, _ = io.WriteString(w, `{"opcionales":[]}`)
		case "/api/v2/arca/tipos-documento":
			_, _ = io.WriteString(w, `{"tipos_documento":[]}`)
		case "/api/v2/arca/monedas":
			_, _ = io.WriteString(w, `{"monedas":[]}`)
		default:
			_, _ = io.WriteString(w, `{"cotizacion":{}}`)
		}
	}))
	defer servidor.Close()

	cliente := nuevoClientePrueba(t, servidor.URL, servidor.Client())
	ctx := context.Background()
	_, _ = cliente.Catalogos.Estado(ctx)
	_, _ = cliente.Catalogos.ListarPuntosVenta(ctx)
	_, _ = cliente.Catalogos.ListarTributos(ctx)
	_, _ = cliente.Catalogos.ListarConceptos(ctx)
	_, _ = cliente.Catalogos.ListarAlicuotas(ctx)
	_, _ = cliente.Catalogos.ListarTiposComprobante(ctx)
	_, _ = cliente.Catalogos.ListarOpcionales(ctx)
	_, _ = cliente.Catalogos.ListarTiposDocumento(ctx)
	_, _ = cliente.Catalogos.ListarMonedas(ctx)
	_, _ = cliente.Catalogos.ObtenerCotizacion(ctx, "DOL")

	if indice != len(esperadas) {
		t.Fatalf("se realizaron %d llamadas, esperadas %d", indice, len(esperadas))
	}
}

func TestPadronConstruyeRutasSeparadas(t *testing.T) {
	esperadas := []string{
		"/api/v2/arca/padron/p4/personas/30123456789",
		"/api/v2/arca/padron/p5/personas/30123456789",
		"/api/v2/arca/padron/p10/personas/30123456789",
		"/api/v2/arca/padron/p13/personas/dni/12345678",
	}
	indice := 0
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != esperadas[indice] {
			t.Errorf("ruta = %s, esperada %s", r.URL.Path, esperadas[indice])
		}
		indice++
		if r.Header.Get("X-ARCA-Ambiente") != "homologacion" {
			t.Error("falta ambiente ARCA")
		}
		_, _ = io.WriteString(w, `{}`)
	}))
	defer servidor.Close()

	cliente := nuevoClientePrueba(t, servidor.URL, servidor.Client())
	ctx := context.Background()
	_, _ = cliente.Padron.ConsultarP4(ctx, "30123456789")
	_, _ = cliente.Padron.ConsultarP5(ctx, "30123456789")
	_, _ = cliente.Padron.ConsultarP10(ctx, "30123456789")
	_, _ = cliente.Padron.ConsultarP13(ctx, padron.TipoDocumentoDNI, "12345678")
}

func TestPadronP13RechazaTipoDocumentoDesconocido(t *testing.T) {
	cliente, err := NuevoCliente(OpcionesCliente{
		ClaveAPI: "clave-prueba",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cliente.Padron.ConsultarP13(
		context.Background(),
		padron.TipoDocumentoP13("pasaporte"),
		"123",
	); err == nil {
		t.Fatal("se esperaba error para un tipo de documento desconocido")
	}
}

func TestEmitirFacturaEnviaWrapperYCabeceras(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v2/arca/factura" {
			t.Fatalf("solicitud inesperada: %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("Idempotency-Key") != "factura-1" {
			t.Errorf("Idempotency-Key inesperada")
		}
		if r.Header.Get("X-External-Reference") != "pedido-9" {
			t.Errorf("X-External-Reference inesperada")
		}
		if r.Header.Get("X-ARCA-Ambiente") != "homologacion" {
			t.Errorf("ambiente inesperado")
		}

		var cuerpo map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&cuerpo); err != nil {
			t.Fatalf("decodificar cuerpo: %v", err)
		}
		if len(cuerpo) != 1 || cuerpo["data"] == nil {
			t.Fatalf("wrapper inesperado: %#v", cuerpo)
		}
		var data struct {
			Factura struct {
				TipoComprobante int   `json:"tipo_comprobante"`
				Detalles        []any `json:"detalles"`
			} `json:"factura"`
		}
		if err := json.Unmarshal(cuerpo["data"], &data); err != nil {
			t.Fatalf("decodificar data: %v", err)
		}
		if data.Factura.TipoComprobante != 11 || len(data.Factura.Detalles) != 1 {
			t.Fatalf("factura inesperada: %#v", data.Factura)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, `{"documento_id":"doc-1","factura":{"cae":"123"}}`)
	}))
	defer servidor.Close()

	cliente := nuevoClientePrueba(t, servidor.URL, servidor.Client())
	respuesta, err := cliente.Facturacion.EmitirFactura(
		context.Background(),
		facturacion.SolicitudEmitirFactura{
			Factura: facturacion.Comprobante{
				TipoComprobante: facturacion.ComprobanteFacturaC,
				Concepto:        facturacion.ConceptoProductos,
				TipoDocumento:   facturacion.DocumentoConsumidorFinal,
				CondicionFiscal: facturacion.CondicionConsumidorFinal,
				Detalles: []facturacion.ItemComprobante{{
					ID:             "1",
					Descripcion:    "Servicio",
					PrecioUnitario: 100,
					Cantidad:       1,
					AlicuotasIVA:   facturacion.AlicuotaCero,
					Subtotal:       100,
				}},
			},
		},
		OpcionesEmisionFactura{
			ClaveIdempotencia: "factura-1",
			ReferenciaExterna: "pedido-9",
		},
	)
	if err != nil {
		t.Fatalf("EmitirFactura: %v", err)
	}
	if respuesta.DocumentoID == nil || *respuesta.DocumentoID != "doc-1" {
		t.Fatalf("respuesta inesperada: %#v", respuesta)
	}
}

func TestErroresHTTPSeDevuelvenTipados(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(http.StatusConflict)
		_, _ = io.WriteString(w, `{
			"code":"idempotency_conflict",
			"type":"about:blank",
			"title":"Conflicto",
			"status":409,
			"detail":"La clave ya fue usada",
			"instance":"/api/v2/arca/factura",
			"request_id":"req-1"
		}`)
	}))
	defer servidor.Close()

	cliente := nuevoClientePrueba(t, servidor.URL, servidor.Client())
	_, err := cliente.Catalogos.Estado(context.Background())
	var apiErr *ErrorAPI
	if !errors.As(err, &apiErr) {
		t.Fatalf("se esperaba ErrorAPI, se recibió %T: %v", err, err)
	}
	if apiErr.Estado != 409 || apiErr.Codigo != "idempotency_conflict" || apiErr.IDSolicitud != "req-1" {
		t.Fatalf("ErrorAPI inesperado: %#v", apiErr)
	}
}

func TestJSONInvalidoYContextoCanceladoNoProducenPanic(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{`)
	}))
	defer servidor.Close()

	cliente := nuevoClientePrueba(t, servidor.URL, servidor.Client())
	_, err := cliente.Catalogos.Estado(context.Background())
	if err == nil || !strings.Contains(err.Error(), "decodificar respuesta JSON") {
		t.Fatalf("error inesperado: %v", err)
	}

	ctx, cancelar := context.WithCancel(context.Background())
	cancelar()
	_, err = cliente.Catalogos.Estado(ctx)
	if err == nil {
		t.Fatal("se esperaba error por contexto cancelado")
	}
}
