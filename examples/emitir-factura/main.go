package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	gofip "github.com/unkcode-org/gofip-sdk-go/v2"
	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
)

func main() {
	cliente, err := gofip.NuevoCliente(gofip.OpcionesCliente{
		ClaveAPI:     os.Getenv("GOFIP_API_KEY"),
		URLBase:      os.Getenv("GOFIP_API_URL"),
		AmbienteARCA: gofip.AmbienteHomologacion,
	})
	if err != nil {
		log.Fatal("no se pudo completar la operación GOFIP")
	}

	clave := os.Getenv("GOFIP_IDEMPOTENCY_KEY")
	if clave == "" {
		log.Fatal("falta GOFIP_IDEMPOTENCY_KEY estable por venta")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	respuesta, err := cliente.Facturacion.EmitirFactura(
		ctx,
		facturacion.SolicitudEmitirFactura{
			Factura: facturacion.Comprobante{
				TipoComprobante: facturacion.ComprobanteFacturaC,
				Concepto:        facturacion.ConceptoProductos,
				TipoDocumento:   facturacion.DocumentoConsumidorFinal,
				CondicionFiscal: facturacion.CondicionConsumidorFinal,
				NumeroDocumento: 0,
				Detalles: []facturacion.ItemComprobante{{
					ID:             "SERV-222",
					Descripcion:    "Diseño y desarrollo de sitio web",
					PrecioUnitario: 200000,
					Cantidad:       1,
					AlicuotasIVA:   facturacion.AlicuotaCero,
					Subtotal:       200000,
				}},
			},
		},
		gofip.OpcionesEmisionFactura{
			ClaveIdempotencia: clave,
			ReferenciaExterna: os.Getenv("GOFIP_EXTERNAL_REFERENCE"),
		},
	)
	if err != nil {
		log.Fatal("no se pudo completar la operación GOFIP")
	}

	if respuesta.DocumentoID != nil {
		fmt.Println("Documento:", *respuesta.DocumentoID)
	}
	if respuesta.Factura.CAE != nil {
		fmt.Println("CAE:", *respuesta.Factura.CAE)
	}
}
