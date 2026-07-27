package gofip

import (
	"net/http"
	"testing"
	"time"
)

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
