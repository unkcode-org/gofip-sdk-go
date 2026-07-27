package main

import (
	"context"
	"fmt"
	"log"
	"os"

	gofip "github.com/unkcode-org/gofip-sdk-go/v2"
)

func main() {
	cliente, err := gofip.NuevoCliente(gofip.OpcionesCliente{
		ClaveAPI: os.Getenv("GOFIP_API_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}

	monedas, err := cliente.Catalogos.ListarMonedas(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", monedas.Monedas)
}
