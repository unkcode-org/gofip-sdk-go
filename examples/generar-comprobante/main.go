package main

import (
	"context"
	"log"
	"os"

	gofip "github.com/unkcode-org/gofip-sdk-go/v2"
	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/catalogos"
	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/documentos"
	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
)

func main() {
	cliente, err := gofip.NuevoCliente(gofip.OpcionesCliente{
		ClaveAPI: os.Getenv("GOFIP_API_KEY"),
	})
	if err != nil {
		log.Fatal("no se pudo completar la operación de documentos")
	}

	err = cliente.Documentos.GenerarYGuardarComprobante(
		context.Background(),
		documentos.SolicitudDocumentoComprobante{
			Comprobante: facturacion.ComprobanteAutorizado{},
			Metadatos: facturacion.MetadatosComprobante{
				CondicionFiscalCliente: facturacion.CondicionConsumidorFinal,
				MetodoPago:             catalogos.MedioPagoEfectivo,
			},
			Items: []facturacion.ItemComprobante{},
		},
		"documentos/comprobante.pdf",
		gofip.OpcionesGuardado{Sobrescribir: false},
	)
	if err != nil {
		log.Fatal("no se pudo completar la operación de documentos")
	}
}
