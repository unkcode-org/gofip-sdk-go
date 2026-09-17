package gofip

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripperPrueba func(*http.Request) (*http.Response, error)

func (f roundTripperPrueba) RoundTrip(solicitud *http.Request) (*http.Response, error) {
	return f(solicitud)
}

func TestNuevoClienteAplicaValoresPredeterminados(t *testing.T) {
	cliente, err := NuevoCliente(OpcionesCliente{ClaveAPI: "secreta"})
	if err != nil {
		t.Fatalf("NuevoCliente: %v", err)
	}

	if cliente.Facturacion == nil || cliente.Documentos == nil || cliente.Padron == nil ||
		cliente.Catalogos == nil || cliente.Utilidades == nil {
		t.Fatal("se esperaban todos los servicios inicializados")
	}
	if cliente.Facturacion.transporte.urlBase != URLBasePredeterminada {
		t.Fatalf("URL base inesperada: %s", cliente.Facturacion.transporte.urlBase)
	}
	if cliente.Facturacion.transporte.ambienteARCA != AmbienteProduccion {
		t.Fatalf("ambiente inesperado: %s", cliente.Facturacion.transporte.ambienteARCA)
	}
	if cliente.Facturacion.transporte.cliente.Timeout != 90*time.Second {
		t.Fatalf("timeout inesperado: %s", cliente.Facturacion.transporte.cliente.Timeout)
	}
}

func TestNuevoClienteValidaConfiguracion(t *testing.T) {
	pruebas := []struct {
		nombre   string
		opciones OpcionesCliente
	}{
		{"clave vacia", OpcionesCliente{}},
		{"URL relativa", OpcionesCliente{ClaveAPI: "key", URLBase: "/api"}},
		{"esquema invalido", OpcionesCliente{ClaveAPI: "key", URLBase: "ftp://example.com"}},
		{"ambiente invalido", OpcionesCliente{ClaveAPI: "key", AmbienteARCA: "pruebas"}},
	}

	for _, prueba := range pruebas {
		t.Run(prueba.nombre, func(t *testing.T) {
			if _, err := NuevoCliente(prueba.opciones); err == nil {
				t.Fatal("se esperaba error")
			}
		})
	}
}

func TestNuevoClienteUsaClienteHTTPPersonalizado(t *testing.T) {
	var solicitudRecibida *http.Request
	clienteHTTP := &http.Client{Transport: roundTripperPrueba(func(solicitud *http.Request) (*http.Response, error) {
		solicitudRecibida = solicitud
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader("disponible")),
			Request:    solicitud,
		}, nil
	})}

	cliente, err := NuevoCliente(OpcionesCliente{
		ClaveAPI:    "clave-prueba",
		URLBase:     "https://api.example.test",
		ClienteHTTP: clienteHTTP,
	})
	if err != nil {
		t.Fatalf("NuevoCliente: %v", err)
	}
	if cliente.Utilidades.transporte.cliente != clienteHTTP {
		t.Fatal("el cliente HTTP personalizado no se conservó")
	}

	respuesta, err := cliente.Utilidades.Disponibilidad(context.Background())
	if err != nil {
		t.Fatalf("Disponibilidad: %v", err)
	}
	if respuesta != "disponible" {
		t.Fatalf("respuesta = %q, esperada disponible", respuesta)
	}
	if solicitudRecibida == nil {
		t.Fatal("el transporte personalizado no recibió una solicitud")
	}
	if solicitudRecibida.URL.String() != "https://api.example.test/api/v2/disponibilidad" {
		t.Fatalf("URL = %s", solicitudRecibida.URL)
	}
	if solicitudRecibida.Header.Get("Authorization") != "Bearer clave-prueba" {
		t.Fatal("la solicitud no conservó la autenticación esperada")
	}

	archivo, err := cliente.Documentos.DescargarPDFRecibo(
		context.Background(),
		"https://archivos.example.test/documento.pdf",
	)
	if err != nil {
		t.Fatalf("DescargarPDFRecibo: %v", err)
	}
	if string(archivo.Contenido) != "disponible" {
		t.Fatalf("contenido = %q, esperado disponible", archivo.Contenido)
	}
	if solicitudRecibida.URL.String() != "https://archivos.example.test/documento.pdf" {
		t.Fatalf("URL de descarga = %s", solicitudRecibida.URL)
	}
}

func nuevoClientePrueba(t *testing.T, urlBase string, clienteHTTP *http.Client) *Cliente {
	t.Helper()
	cliente, err := NuevoCliente(OpcionesCliente{
		ClaveAPI:     "clave-prueba",
		URLBase:      urlBase,
		AmbienteARCA: AmbienteHomologacion,
		ClienteHTTP:  clienteHTTP,
	})
	if err != nil {
		t.Fatalf("NuevoCliente: %v", err)
	}
	return cliente
}
