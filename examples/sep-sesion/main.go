package main

import (
	"context"
	"fmt"
	gofip "github.com/unkcode-org/gofip-sdk-go/v2"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	c, err := gofip.NuevoClienteSEP(gofip.OpcionesClienteSEP{ClaveIntegracion: os.Getenv("GOFIP_SEP_KEY"), URLBase: os.Getenv("GOFIP_API_URL"), ClienteHTTP: &http.Client{Timeout: 15 * time.Second}})
	if err != nil {
		log.Fatal("configuración SEP inválida")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	respuesta, err := c.CrearSesion(ctx)
	if err != nil {
		log.Fatal("no se pudo completar la operación SEP")
	}
	// No imprimir respuesta ni capacidades. El integrador entrega Valor()
	// deliberadamente a su frontend por un canal autorizado, sin cache ni logs.
	fmt.Printf("Operación SEP completada; vencimiento: %s\n", respuesta.Datos.ExpiraDocumentoEn.Format(time.RFC3339))
}
