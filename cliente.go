package gofip

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const URLBasePredeterminada = "https://arca.apps.staging.unkcode.com"

type AmbienteARCA string

const (
	AmbienteProduccion   AmbienteARCA = "produccion"
	AmbienteHomologacion AmbienteARCA = "homologacion"
)

type OpcionesCliente struct {
	ClaveAPI     string
	URLBase      string
	AmbienteARCA AmbienteARCA
	ClienteHTTP  *http.Client
}

type Cliente struct {
	Utilidades  *ServicioUtilidades
	Catalogos   *ServicioCatalogos
	Padron      *ServicioPadron
	Facturacion *ServicioFacturacion
	Documentos  *ServicioDocumentos
}

func NuevoCliente(opciones OpcionesCliente) (*Cliente, error) {
	opciones.ClaveAPI = strings.TrimSpace(opciones.ClaveAPI)
	if opciones.ClaveAPI == "" {
		return nil, errors.New("la clave API es obligatoria")
	}

	if opciones.URLBase == "" {
		opciones.URLBase = URLBasePredeterminada
	}
	opciones.URLBase = strings.TrimRight(opciones.URLBase, "/")

	u, err := url.Parse(opciones.URLBase)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, errors.New("la URL base debe ser una URL HTTP o HTTPS absoluta")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("la URL base debe usar HTTP o HTTPS")
	}

	if opciones.AmbienteARCA == "" {
		opciones.AmbienteARCA = AmbienteProduccion
	}
	if opciones.AmbienteARCA != AmbienteProduccion && opciones.AmbienteARCA != AmbienteHomologacion {
		return nil, errors.New("el ambiente ARCA debe ser produccion u homologacion")
	}

	if opciones.ClienteHTTP == nil {
		opciones.ClienteHTTP = &http.Client{Timeout: 90 * time.Second}
	}

	t := &transporte{
		claveAPI:     opciones.ClaveAPI,
		urlBase:      opciones.URLBase,
		ambienteARCA: opciones.AmbienteARCA,
		cliente:      opciones.ClienteHTTP,
	}

	return &Cliente{
		Utilidades:  &ServicioUtilidades{transporte: t},
		Catalogos:   &ServicioCatalogos{transporte: t},
		Padron:      &ServicioPadron{transporte: t},
		Facturacion: &ServicioFacturacion{transporte: t},
		Documentos:  &ServicioDocumentos{transporte: t},
	}, nil
}
