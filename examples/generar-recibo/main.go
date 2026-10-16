package main

import (
	"context"
	"fmt"
	"log"
	"os"

	gofip "github.com/unkcode-org/gofip-sdk-go/v2"
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

	ctx := context.Background()
	respuesta, err := cliente.Documentos.GenerarRecibo(
		ctx,
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
			},
			Items: []documentos.ItemRecibo{{
				ID:          "SERV-001",
				Descripcion: "Servicio profesional",
				Cantidad:    1,
				Precio:      12500.50,
				Total:       12500.50,
			}},
		},
		gofip.OpcionesGenerarRecibo{ClaveIdempotencia: "recibo-pago-123"},
	)
	if err != nil {
		log.Fatal("no se pudo completar la operación de documentos")
	}

	fmt.Println("Documento:", respuesta.DocumentoID)
	if err := cliente.Documentos.DescargarYGuardarPDFRecibo(
		ctx,
		respuesta.PDFURL,
		"documentos/recibo.pdf",
		gofip.OpcionesGuardado{Sobrescribir: false},
	); err != nil {
		log.Fatal("no se pudo completar la operación de documentos")
	}
}
