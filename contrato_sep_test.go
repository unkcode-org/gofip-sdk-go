package gofip

import (
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/sep"
)

// Este inventario compara las constantes públicas con el enum del YAML canónico
// copiado del backend; evita que una modificación agregue permisos por suposición.
func TestPermisosSEPCoincidenConOpenAPI(t *testing.T) {
	b, err := os.ReadFile("docs/openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	inicio := strings.Index(s, "        PermisoEmbedded:\n")
	if inicio < 0 {
		t.Fatal("schema de permisos ausente")
	}
	bloque := s[inicio+len("        PermisoEmbedded:\n"):]
	if fin := regexp.MustCompile(`(?m)^        [A-Za-z0-9_]+:`).FindStringIndex(bloque); fin != nil {
		bloque = bloque[:fin[0]]
	}
	var publicados []string
	for _, m := range regexp.MustCompile(`"([a-z_]+:[a-z]+)"`).FindAllStringSubmatch(bloque, -1) {
		publicados = append(publicados, m[1])
	}
	constantes := []sep.PermisoEmbedded{sep.PermisoFacturasLeer, sep.PermisoFacturasCrear, sep.PermisoClientesLeer, sep.PermisoClientesEscribir, sep.PermisoPuntosVentaLeer, sep.PermisoPuntosVentaEscribir, sep.PermisoRecurrentesLeer, sep.PermisoRecurrentesEscribir, sep.PermisoWhatsAppLeer, sep.PermisoWhatsAppEscribir, sep.PermisoARCALeer, sep.PermisoARCAEscribir}
	var implementados []string
	for _, p := range constantes {
		implementados = append(implementados, string(p))
	}
	sort.Strings(publicados)
	sort.Strings(implementados)
	if len(publicados) == 0 || strings.Join(publicados, ",") != strings.Join(implementados, ",") {
		t.Fatalf("permisos divergentes del contrato: %v / %v", publicados, implementados)
	}
}
