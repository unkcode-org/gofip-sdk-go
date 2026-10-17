package gofip

import (
	"crypto/sha256"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
)

func TestContratoImportacionYProcedenciaExacta(t *testing.T) {
	yaml, err := os.ReadFile("docs/openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}
	provenance, err := os.ReadFile("docs/CONTRATO.md")
	if err != nil {
		t.Fatal(err)
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(yaml))
	if !strings.Contains(string(provenance), hash) || !strings.Contains(string(provenance), "Versión: 2.19.0.") {
		t.Fatal("copia OpenAPI sin procedencia o checksum actualizado")
	}
	for _, token := range []string{`version: "2.19.0"`, "/api/v2/facturas/importaciones:", "facturas:import", "legacy_import", "SolicitudImportarFactura:", "RespuestaImportarFactura:", "DecimalHistorico:", "VerificacionImportacion:", "CompletitudDocumental:"} {
		if !strings.Contains(string(yaml), token) {
			t.Fatalf("contrato sin %s", token)
		}
	}
	for _, typ := range []reflect.Type{reflect.TypeOf(OpcionesImportacionFactura{}), reflect.TypeOf(facturacion.SolicitudImportarFactura{}), reflect.TypeOf(facturacion.FacturaHistorica{})} {
		for i := 0; i < typ.NumField(); i++ {
			if regexp.MustCompile(`(?i)ambiente|emisor`).MatchString(typ.Field(i).Name) {
				t.Fatalf("%s amplía autoridad", typ)
			}
		}
	}
}
