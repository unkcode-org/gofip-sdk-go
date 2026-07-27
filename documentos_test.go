package gofip

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/documentos"
	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
)

func TestGenerarReciboEnviaContratoEspecificoSinWrapperData(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/templates-oficiales/recibo/documentos" {
			t.Fatalf("ruta inesperada: %s", r.URL.Path)
		}
		if r.Header.Get("Idempotency-Key") != "recibo-1" {
			t.Errorf("Idempotency-Key inesperada")
		}
		if r.Header.Get("X-ARCA-Ambiente") != "" {
			t.Errorf("el endpoint de templates no debe recibir X-ARCA-Ambiente")
		}

		var cuerpo map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&cuerpo); err != nil {
			t.Fatalf("decodificar cuerpo: %v", err)
		}
		if len(cuerpo) != 1 || cuerpo["datos"] == nil || cuerpo["data"] != nil {
			t.Fatalf("cuerpo exterior inesperado: %#v", cuerpo)
		}
		var datos documentos.DatosRecibo
		if err := json.Unmarshal(cuerpo["datos"], &datos); err != nil {
			t.Fatalf("decodificar datos: %v", err)
		}
		if datos.Recibo.NumeroRecibo != "REC-000001" || len(datos.Items) != 1 {
			t.Fatalf("datos inesperados: %#v", datos)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = io.WriteString(w, `{
			"documento_id":"doc-recibo",
			"template_slug":"recibo",
			"template_version":"1.0.0",
			"pdf_url":"https://files.example/recibo.pdf?token=abc",
			"pdf_content_type":"application/pdf",
			"pdf_generated_at":"2026-07-27T03:00:00Z"
		}`)
	}))
	defer servidor.Close()

	cliente := nuevoClientePrueba(t, servidor.URL, servidor.Client())
	telefono := "+54 11 5555-1234"
	email := "administracion@cliente.com"
	respuesta, err := cliente.Documentos.GenerarRecibo(
		context.Background(),
		documentos.DatosRecibo{
			Recibo: documentos.Recibo{
				NumeroRecibo: "REC-000001",
				FechaRecibo:  "20260727",
				Subtotal:     12500.50,
				Total:        12500.50,
			},
			Metadatos: documentos.MetadatosRecibo{
				CondicionFiscalCliente:    facturacion.CondicionConsumidorFinal,
				DocumentoCliente:          30123456789,
				NombreORazonSocialCliente: "Cliente de Ejemplo S.A.",
				DomicilioCliente:          "Av. Corrientes 1234, CABA",
				TelefonoCliente:           &telefono,
				EmailCliente:              &email,
			},
			Items: []documentos.ItemRecibo{{
				ID:          "SERV-001",
				Descripcion: "Servicio profesional",
				Cantidad:    1,
				Precio:      12500.50,
				Total:       12500.50,
			}},
		},
		OpcionesGenerarRecibo{ClaveIdempotencia: "recibo-1"},
	)
	if err != nil {
		t.Fatalf("GenerarRecibo: %v", err)
	}
	if respuesta.TemplateSlug != documentos.SlugRecibo ||
		respuesta.TemplateVersion != documentos.VersionRecibo {
		t.Fatalf("respuesta inesperada: %#v", respuesta)
	}
}

func TestGenerarComprobanteDevuelveBinarioYMetadatos(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/documentos/comprobante" {
			t.Fatalf("ruta inesperada: %s", r.URL.Path)
		}
		var cuerpo map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&cuerpo); err != nil {
			t.Fatalf("decodificar cuerpo: %v", err)
		}
		if cuerpo["data"] == nil {
			t.Fatal("falta wrapper data")
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", `attachment; filename="factura.pdf"`)
		_, _ = w.Write([]byte("%PDF-prueba"))
	}))
	defer servidor.Close()

	cliente := nuevoClientePrueba(t, servidor.URL, servidor.Client())
	archivo, err := cliente.Documentos.GenerarComprobante(
		context.Background(),
		documentos.SolicitudDocumentoComprobante{},
	)
	if err != nil {
		t.Fatalf("GenerarComprobante: %v", err)
	}
	if string(archivo.Contenido) != "%PDF-prueba" ||
		archivo.TipoContenido != "application/pdf" ||
		archivo.Nombre != "factura.pdf" {
		t.Fatalf("archivo inesperado: %#v", archivo)
	}
}

func TestDescargaFirmadaNoFiltraCredenciales(t *testing.T) {
	servidorArchivos := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "" {
			t.Errorf("la descarga firmada filtró Authorization: %q", got)
		}
		if got := r.Header.Get("X-ARCA-Ambiente"); got != "" {
			t.Errorf("la descarga firmada filtró ambiente: %q", got)
		}
		w.Header().Set("Content-Type", "application/pdf")
		_, _ = w.Write([]byte("pdf"))
	}))
	defer servidorArchivos.Close()

	cliente := nuevoClientePrueba(t, servidorArchivos.URL, servidorArchivos.Client())
	archivo, err := cliente.Documentos.DescargarPDFRecibo(context.Background(), servidorArchivos.URL+"/signed")
	if err != nil {
		t.Fatalf("DescargarPDFRecibo: %v", err)
	}
	if string(archivo.Contenido) != "pdf" {
		t.Fatalf("contenido inesperado: %q", archivo.Contenido)
	}
}

func TestGuardarArchivoRespetaOpcionSobrescribir(t *testing.T) {
	directorio := t.TempDir()
	ruta := filepath.Join(directorio, "anidado", "recibo.pdf")
	archivo := &Archivo{Contenido: []byte("primero")}

	if err := guardarArchivo(archivo, ruta, OpcionesGuardado{}); err != nil {
		t.Fatalf("primer guardado: %v", err)
	}
	if err := guardarArchivo(&Archivo{Contenido: []byte("segundo")}, ruta, OpcionesGuardado{}); err == nil {
		t.Fatal("se esperaba error al guardar sin sobrescribir")
	}
	contenido, err := os.ReadFile(ruta)
	if err != nil {
		t.Fatalf("leer archivo: %v", err)
	}
	if string(contenido) != "primero" {
		t.Fatalf("el archivo fue modificado: %q", contenido)
	}

	if err := guardarArchivo(
		&Archivo{Contenido: []byte("segundo")},
		ruta,
		OpcionesGuardado{Sobrescribir: true},
	); err != nil {
		t.Fatalf("sobrescribir: %v", err)
	}
	contenido, _ = os.ReadFile(ruta)
	if string(contenido) != "segundo" {
		t.Fatalf("contenido sobrescrito inesperado: %q", contenido)
	}
}

func TestGuardarArchivoLimpiaParcialAnteError(t *testing.T) {
	directorio := t.TempDir()
	rutaDirectorio := filepath.Join(directorio, "destino")
	if err := os.Mkdir(rutaDirectorio, 0o755); err != nil {
		t.Fatal(err)
	}
	err := guardarArchivo(&Archivo{Contenido: []byte("contenido")}, rutaDirectorio, OpcionesGuardado{
		Sobrescribir: true,
	})
	if err == nil {
		t.Fatal("se esperaba error al usar un directorio como archivo")
	}
}

func TestListarTemplatesDecodificaTimestamp(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = io.WriteString(w, `{"templates":[{
			"slug":"recibo",
			"version":"1.0.0",
			"nombre":"Recibo",
			"descripcion":"Comprobante de pago",
			"autor":"GOFIP",
			"ultima_modificacion":"2026-07-27T03:00:00Z",
			"schema_datos":{"type":"object"}
		}]}`)
	}))
	defer servidor.Close()

	cliente := nuevoClientePrueba(t, servidor.URL, servidor.Client())
	respuesta, err := cliente.Documentos.ListarTemplatesOficiales(context.Background())
	if err != nil {
		t.Fatalf("ListarTemplatesOficiales: %v", err)
	}
	if len(respuesta.Templates) != 1 || respuesta.Templates[0].UltimaModificacion.IsZero() {
		t.Fatalf("templates inesperados: %#v", respuesta.Templates)
	}
}

func TestErrorDescargaNoJSON(t *testing.T) {
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = io.WriteString(w, "almacenamiento no disponible")
	}))
	defer servidor.Close()

	cliente := nuevoClientePrueba(t, servidor.URL, servidor.Client())
	_, err := cliente.Documentos.DescargarPDFRecibo(context.Background(), servidor.URL)
	var httpErr *ErrorHTTP
	if !errors.As(err, &httpErr) || httpErr.Estado != http.StatusBadGateway {
		t.Fatalf("error inesperado: %T %v", err, err)
	}
}
