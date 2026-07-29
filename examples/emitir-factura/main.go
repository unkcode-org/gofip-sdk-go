package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	gofip "github.com/unkcode-org/gofip-sdk-go/v2"
	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
)

func main() {
	cliente, err := gofip.NuevoCliente(gofip.OpcionesCliente{
		ClaveAPI:     os.Getenv("GOFIP_API_KEY"),
		AmbienteARCA: gofip.AmbienteHomologacion,
	})
	if err != nil {
		log.Fatal(err)
	}

	respuesta, err := cliente.Facturacion.EmitirFactura(
		context.Background(),
		facturacion.SolicitudEmitirFactura{
			Factura: facturacion.Comprobante{
				TipoComprobante: facturacion.ComprobanteFacturaC,
				Concepto:        facturacion.ConceptoProductos,
				TipoDocumento:   facturacion.DocumentoConsumidorFinal,
				CondicionFiscal: facturacion.CondicionConsumidorFinal,
				Detalles: []facturacion.ItemComprobante{{
					ID:             "SERV-001",
					Descripcion:    "Servicio profesional",
					PrecioUnitario: 1200.50,
					Cantidad:       1,
					AlicuotasIVA:   facturacion.AlicuotaCero,
					Subtotal:       1200.50,
				}},
			},
		},
		gofip.OpcionesEmisionFactura{
			ClaveIdempotencia: "factura-pedido-125",
			ReferenciaExterna: "pedido-125",
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	json, err := json.Marshal(respuesta)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n", string(json))
	fmt.Printf("%+v\n", respuesta)
}
