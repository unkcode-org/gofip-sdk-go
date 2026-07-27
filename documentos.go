package gofip

import (
	"context"
	"net/http"
	"net/url"

	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/documentos"
)

type ServicioDocumentos struct {
	transporte *transporte
}

type OpcionesGenerarRecibo struct {
	ClaveIdempotencia string
}

func (s *ServicioDocumentos) GenerarComprobante(
	ctx context.Context,
	solicitud documentos.SolicitudDocumentoComprobante,
) (*Archivo, error) {
	return s.transporte.binario(
		ctx,
		http.MethodPost,
		"/api/v2/documentos/comprobante",
		cuerpoDatos[documentos.SolicitudDocumentoComprobante]{Datos: solicitud},
		nil,
	)
}

func (s *ServicioDocumentos) GenerarYGuardarComprobante(
	ctx context.Context,
	solicitud documentos.SolicitudDocumentoComprobante,
	ruta string,
	opciones OpcionesGuardado,
) error {
	archivo, err := s.GenerarComprobante(ctx, solicitud)
	if err != nil {
		return err
	}
	return guardarArchivo(archivo, ruta, opciones)
}

func (s *ServicioDocumentos) ListarTemplatesOficiales(
	ctx context.Context,
) (*documentos.RespuestaTemplatesOficiales, error) {
	var respuesta documentos.RespuestaTemplatesOficiales
	if err := s.transporte.json(
		ctx,
		http.MethodGet,
		"/api/v2/templates-oficiales",
		nil,
		nil,
		false,
		&respuesta,
	); err != nil {
		return nil, err
	}
	return &respuesta, nil
}

func (s *ServicioDocumentos) GenerarRecibo(
	ctx context.Context,
	datos documentos.DatosRecibo,
	opciones OpcionesGenerarRecibo,
) (*documentos.RespuestaDocumentoNoFiscal, error) {
	var respuesta documentos.RespuestaDocumentoNoFiscal
	ruta := "/api/v2/templates-oficiales/" + url.PathEscape(documentos.SlugRecibo) + "/documentos"
	cabeceras := map[string]string{"Idempotency-Key": opciones.ClaveIdempotencia}
	cuerpo := documentos.SolicitudGenerarRecibo{Datos: datos}
	if err := s.transporte.json(ctx, http.MethodPost, ruta, cuerpo, cabeceras, false, &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}

func (s *ServicioDocumentos) ObtenerPDFRecibo(
	ctx context.Context,
	documentoID string,
) (*documentos.RespuestaDocumentoNoFiscal, error) {
	var respuesta documentos.RespuestaDocumentoNoFiscal
	ruta := "/api/v2/documentos-no-fiscales/" + url.PathEscape(documentoID) + "/pdf"
	if err := s.transporte.json(ctx, http.MethodGet, ruta, nil, nil, false, &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}

func (s *ServicioDocumentos) DescargarPDFRecibo(
	ctx context.Context,
	urlFirmada string,
) (*Archivo, error) {
	return s.transporte.descargar(ctx, urlFirmada)
}

func (s *ServicioDocumentos) DescargarYGuardarPDFRecibo(
	ctx context.Context,
	urlFirmada string,
	ruta string,
	opciones OpcionesGuardado,
) error {
	archivo, err := s.DescargarPDFRecibo(ctx, urlFirmada)
	if err != nil {
		return err
	}
	return guardarArchivo(archivo, ruta, opciones)
}

func (s *ServicioDocumentos) GenerarPDFFactura(
	ctx context.Context,
	documentoID string,
	solicitud documentos.SolicitudGenerarPDFFactura,
) (*documentos.RespuestaPDFFactura, error) {
	var respuesta documentos.RespuestaPDFFactura
	ruta := "/api/v2/factura/" + url.PathEscape(documentoID) + "/pdf"
	cuerpo := cuerpoDatos[documentos.SolicitudGenerarPDFFactura]{Datos: solicitud}
	if err := s.transporte.json(ctx, http.MethodPost, ruta, cuerpo, nil, false, &respuesta); err != nil {
		return nil, err
	}
	return &respuesta, nil
}

func (s *ServicioDocumentos) DescargarPDFFactura(
	ctx context.Context,
	urlFirmada string,
) (*Archivo, error) {
	return s.transporte.descargar(ctx, urlFirmada)
}

func (s *ServicioDocumentos) DescargarYGuardarPDFFactura(
	ctx context.Context,
	urlFirmada string,
	ruta string,
	opciones OpcionesGuardado,
) error {
	archivo, err := s.DescargarPDFFactura(ctx, urlFirmada)
	if err != nil {
		return err
	}
	return guardarArchivo(archivo, ruta, opciones)
}
