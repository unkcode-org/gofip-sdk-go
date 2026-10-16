package main

import (
	"context"
	"errors"
	"fmt"
	gofip "github.com/unkcode-org/gofip-sdk-go/v2"
	"log"
	"os"
	"time"
)

func main() {
	c, err := gofip.NuevoCliente(gofip.OpcionesCliente{ClaveAPI: os.Getenv("GOFIP_API_KEY"), URLBase: os.Getenv("GOFIP_API_URL")})
	if err != nil {
		log.Fatal("configuración GOFIP inválida")
	}
	valor := os.Getenv("GOFIP_CAE")
	if valor == "" {
		log.Fatal("falta GOFIP_CAE")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	f, err := c.Facturacion.ObtenerFacturaPorCAE(ctx, valor)
	if err != nil {
		var api *gofip.ErrorAPI
		if errors.As(err, &api) {
			log.Fatalf("consulta fiscal: HTTP %d código %s", api.Estado, api.Codigo)
		}
		log.Fatal("no se pudo consultar la factura")
	}
	fmt.Printf("Documento: %s; estado: %s; CAE: %s\n", f.DocumentoID, f.Estado, f.Factura.CAE)
}
