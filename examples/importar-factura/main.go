// Importa una única autorización Legacy existente. Ejecutar sólo con una
// solicitud revisada y la API key normal del mismo emisor/PV en el entorno autorizado.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	gofip "github.com/unkcode-org/gofip-sdk-go/v2"
	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
)

func main() {
	path, key := os.Getenv("GOFIP_IMPORT_REQUEST_FILE"), os.Getenv("GOFIP_IMPORT_IDEMPOTENCY_KEY")
	if path == "" || key == "" {
		log.Fatal("faltan GOFIP_IMPORT_REQUEST_FILE/GOFIP_IMPORT_IDEMPOTENCY_KEY")
	}
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("no se pudo abrir la solicitud histórica")
	}
	defer f.Close()
	var solicitud facturacion.SolicitudImportarFactura
	dec := json.NewDecoder(io.LimitReader(f, (1<<20)+1))
	dec.DisallowUnknownFields()
	if dec.Decode(&solicitud) != nil || dec.Decode(&struct{}{}) != io.EOF {
		log.Fatal("solicitud histórica inválida; usar el DTO sin envelope data")
	}
	c, err := gofip.NuevoCliente(gofip.OpcionesCliente{URLBase: os.Getenv("GOFIP_API_URL"), ClaveAPI: os.Getenv("GOFIP_API_KEY")})
	if err != nil {
		log.Fatal("configuración GOFIP inválida")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	r, err := c.Facturacion.ImportarFactura(ctx, solicitud, gofip.OpcionesImportacionFactura{ClaveIdempotencia: key, ReferenciaExterna: os.Getenv("GOFIP_IMPORT_EXTERNAL_REFERENCE")})
	if err != nil {
		var api *gofip.ErrorAPI
		if errors.As(err, &api) {
			log.Fatalf("importación fiscal: HTTP %d código %s", api.Estado, api.Codigo)
		}
		log.Fatal("no se pudo importar; conservar la clave para reanudar y revisar el resultado")
	}
	fmt.Printf("Documento: %s; importado: %t; ya existía: %t\n", r.DocumentoID, r.Importado, r.YaExistia)
}
