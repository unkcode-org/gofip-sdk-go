package gofip

import (
	"context"
	"errors"
	"net/url"
	"regexp"

	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/sep"
)

var patronBootstrapSEP = regexp.MustCompile(`^gofip_emb_[A-Za-z0-9_-]{43}$`)

// CrearSesion hereda todos los permisos de la Integration Key, sin parámetros de autoridad.
// Entrega un documento de framing y un bootstrap de un solo uso.
// La sesión del navegador sólo nace durante el canje que realiza el frontend.
func (c *ClienteSEP) CrearSesion(ctx context.Context) (*sep.RespuestaCrearSesionEmbedded, error) {
	var r sep.RespuestaCrearSesionEmbedded
	if err := c.post(ctx, "/api/v2/embedded/sessions", &r); err != nil {
		return nil, err
	}
	u, err := url.Parse(r.Datos.URLDocumento.Valor())
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || r.Datos.ExpiraDocumentoEn.IsZero() || r.Datos.ExpiraBootstrapEn.IsZero() || !patronBootstrapSEP.MatchString(r.Datos.TokenBootstrap.Valor()) {
		return nil, errors.New("respuesta de creación SEP inválida")
	}
	return &r, nil
}

// CrearBootstrap emite una nueva capacidad de recuperación sin crear documento
// ni recuperar secretos anteriores. No reintenta automáticamente el POST.
func (c *ClienteSEP) CrearBootstrap(ctx context.Context) (*sep.RespuestaCrearBootstrapEmbedded, error) {
	var r sep.RespuestaCrearBootstrapEmbedded
	if err := c.post(ctx, "/api/v2/embedded/bootstraps", &r); err != nil {
		return nil, err
	}
	if r.Datos.ExpiraBootstrapEn.IsZero() || !patronBootstrapSEP.MatchString(r.Datos.TokenBootstrap.Valor()) {
		return nil, errors.New("respuesta de recuperación SEP inválida")
	}
	return &r, nil
}
