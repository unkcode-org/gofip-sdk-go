package gofip

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/sep"
)

var patronCodigoSEP = regexp.MustCompile(`^[0-9]{5}$`)

func (c *ClienteSEP) post(ctx context.Context, ruta string, solicitud sep.SolicitudCrearSesionEmbedded, destino any) error {
	contenido, err := json.Marshal(cuerpoDatos[sep.SolicitudCrearSesionEmbedded]{Datos: solicitud})
	if err != nil {
		return errors.New("serializar solicitud SEP")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.urlBase+ruta, bytes.NewReader(contenido))
	if err != nil {
		return errors.New("crear solicitud SEP")
	}
	req.Header.Set("Authorization", "IntegrationKey "+c.clave)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	// Sin Idempotency-Key y sin política de reintentos: estos POST crean capacidades.
	resp, err := c.cliente.Do(req)
	if err != nil {
		if seguro := errorContextoSEP(ctx, err); seguro != nil {
			return seguro
		}
		// Un transporte externo puede incluir headers o URLs sensibles en su error.
		return errors.New("realizar solicitud SEP: fallo de transporte")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		// Compartimos el parser existente, pero nunca retenemos campos ni cuerpos
		// remotos arbitrarios en errores SEP: pueden reflejar credenciales/capacidades.
		limitada := *resp
		limitada.Body = io.NopCloser(io.LimitReader(resp.Body, 1<<20))
		original := errorDesdeRespuesta(&limitada)
		if seguro := errorContextoSEP(ctx, original); seguro != nil {
			return seguro
		}
		var api *ErrorAPI
		if errors.As(original, &api) {
			codigo := ""
			if patronCodigoSEP.MatchString(api.Codigo) {
				codigo = api.Codigo
			}
			return &ErrorAPI{Estado: resp.StatusCode, Codigo: codigo, Titulo: "Error de API SEP", Detalle: http.StatusText(resp.StatusCode)}
		}
		return &ErrorHTTP{Estado: resp.StatusCode, Cuerpo: http.StatusText(resp.StatusCode)}
	}
	// Limitar la respuesta y exigir un único objeto evita retener cuerpos grandes
	// o aceptar JSON adicional. Nunca adjuntar el error de decode con datos remotos.
	dec := json.NewDecoder(io.LimitReader(resp.Body, 1<<20))
	if err := dec.Decode(destino); err != nil {
		if seguro := errorContextoSEP(ctx, err); seguro != nil {
			return seguro
		}
		return errors.New("decodificar respuesta SEP inválida")
	}
	if err := dec.Decode(&struct{}{}); err != io.EOF {
		if seguro := errorContextoSEP(ctx, err); seguro != nil {
			return seguro
		}
		return errors.New("respuesta SEP contiene datos adicionales")
	}
	return nil
}

// Preservar sólo causas conocidas y seguras; los errores arbitrarios de un
// transporte o decoder pueden incluir secretos en su texto o cadena Unwrap.
func errorContextoSEP(ctx context.Context, err error) error {
	if errors.Is(err, context.Canceled) || errors.Is(ctx.Err(), context.Canceled) {
		return fmt.Errorf("solicitud SEP: %w", context.Canceled)
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("solicitud SEP: %w", context.DeadlineExceeded)
	}
	return nil
}
