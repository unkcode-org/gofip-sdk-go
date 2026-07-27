package main

import (
	"context"
	"fmt"
	"log"
	"os"

	gofip "github.com/unkcode-org/gofip-sdk-go/v2"
	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/padron"
)

func main() {
	cliente, err := gofip.NuevoCliente(gofip.OpcionesCliente{
		ClaveAPI: os.Getenv("GOFIP_API_KEY"),
	})
	if err != nil {
		log.Fatal(err)
	}

	persona, err := cliente.Padron.ConsultarP13(
		context.Background(),
		padron.TipoDocumentoDNI,
		"12345678",
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", persona.PersonaP13)
}
