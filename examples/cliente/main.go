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

	salud, err := cliente.Utilidades.Salud(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(salud)
}
