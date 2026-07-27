package gofip

import (
	"context"
	"net/http"

	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/catalogos"
)

type ServicioUtilidades struct {
	transporte *transporte
}

func (s *ServicioUtilidades) Salud(ctx context.Context) (string, error) {
	return s.transporte.texto(ctx, "/api/v2/salud", false)
}

func (s *ServicioUtilidades) Disponibilidad(ctx context.Context) (string, error) {
	return s.transporte.texto(ctx, "/api/v2/disponibilidad", true)
}

func (s *ServicioUtilidades) ListarMediosPago(ctx context.Context) (*catalogos.RespuestaMediosPago, error) {
	var respuesta catalogos.RespuestaMediosPago
	if err := s.transporte.json(
		ctx,
		http.MethodGet,
		"/api/v2/medios-pago",
		nil,
		nil,
		false,
		&respuesta,
	); err != nil {
		return nil, err
	}
	return &respuesta, nil
}
