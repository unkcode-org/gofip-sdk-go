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

	nombreCliente := "LERVU S.A."

	err = cliente.Documentos.GenerarYGuardarPDFFactura(
		context.Background(),
		os.Getenv("GOFIP_DOCUMENTO_ID"),
		documentos.SolicitudGenerarPDFFactura{
			Metadatos: facturacion.MetadatosComprobante{
				CondicionFiscalCliente:    facturacion.CondicionConsumidorFinal,
				MetodoPago:                catalogos.MedioPagoEfectivo,
				NombreORazonSocialCliente: &nombreCliente,
			},
		},
		"documentos/factura.pdf",
		gofip.OpcionesGuardado{Sobrescribir: true},
	)
	if err != nil {
		log.Fatal("no se pudo completar la operación de documentos")
	}
}
