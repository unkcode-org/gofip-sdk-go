// Package sep modela exclusivamente las operaciones M2M de GOFIP Embedded.
package sep

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type PermisoEmbedded string

const (
	PermisoFacturasLeer        PermisoEmbedded = "invoices:read"
	PermisoFacturasCrear       PermisoEmbedded = "invoices:create"
	PermisoClientesLeer        PermisoEmbedded = "customers:read"
	PermisoClientesEscribir    PermisoEmbedded = "customers:write"
	PermisoPuntosVentaLeer     PermisoEmbedded = "points_of_sale:read"
	PermisoPuntosVentaEscribir PermisoEmbedded = "points_of_sale:write"
	PermisoRecurrentesLeer     PermisoEmbedded = "recurring_invoices:read"
	PermisoRecurrentesEscribir PermisoEmbedded = "recurring_invoices:write"
	PermisoWhatsAppLeer        PermisoEmbedded = "whatsapp:read"
	PermisoWhatsAppEscribir    PermisoEmbedded = "whatsapp:write"
	PermisoARCALeer            PermisoEmbedded = "arca:read"
	PermisoARCAEscribir        PermisoEmbedded = "arca:write"
)

// Secreto exige acceso explícito mediante Valor. JSON, fmt y slog redactan el
// contenido, también cuando forma parte de una respuesta o un struct anidado.
// JSON de diagnóstico no sirve para reenviar la capacidad al frontend: para
// esa entrega deliberada el integrador debe construir su propio DTO con Valor().
type Secreto struct{ valor string }

func (s Secreto) Valor() string              { return s.valor }
func (Secreto) String() string               { return "[REDACTADO]" }
func (Secreto) GoString() string             { return "[REDACTADO]" }
func (Secreto) Format(f fmt.State, _ rune)   { _, _ = fmt.Fprint(f, "[REDACTADO]") }
func (Secreto) LogValue() slog.Value         { return slog.StringValue("[REDACTADO]") }
func (Secreto) MarshalJSON() ([]byte, error) { return []byte(`"[REDACTADO]"`), nil }
func (s *Secreto) UnmarshalJSON(b []byte) error {
	var v string
	if string(b) == "null" || json.Unmarshal(b, &v) != nil {
		return errors.New("secreto SEP inválido")
	}
	s.valor = v
	return nil
}

type SolicitudCrearSesionEmbedded struct {
	Permisos []PermisoEmbedded `json:"permissions"`
}

// Validar no cambia el orden de los permisos. El techo de
// autorización pertenece al backend y no puede verificarse localmente.
func (s SolicitudCrearSesionEmbedded) Validar() error {
	if len(s.Permisos) == 0 {
		return errors.New("se requiere al menos un permiso Embedded")
	}
	vistos := make(map[PermisoEmbedded]bool, len(s.Permisos))
	for _, p := range s.Permisos {
		switch p {
		case PermisoFacturasLeer, PermisoFacturasCrear, PermisoClientesLeer, PermisoClientesEscribir, PermisoPuntosVentaLeer, PermisoPuntosVentaEscribir, PermisoRecurrentesLeer, PermisoRecurrentesEscribir, PermisoWhatsAppLeer, PermisoWhatsAppEscribir, PermisoARCALeer, PermisoARCAEscribir:
		default:
			return errors.New("permiso Embedded desconocido")
		}
		if vistos[p] {
			return errors.New("los permisos Embedded no pueden repetirse")
		}
		vistos[p] = true
	}
	return nil
}

type RespuestaCrearSesionEmbedded struct {
	Datos DatosCrearSesionEmbedded `json:"data"`
}
type DatosCrearSesionEmbedded struct {
	URLDocumento      Secreto   `json:"document_url"`
	ExpiraDocumentoEn time.Time `json:"document_expires_at"`
	TokenBootstrap    Secreto   `json:"bootstrap_token"`
	ExpiraBootstrapEn time.Time `json:"bootstrap_expires_at"`
}
type RespuestaCrearBootstrapEmbedded struct {
	Datos DatosCrearBootstrapEmbedded `json:"data"`
}
type DatosCrearBootstrapEmbedded struct {
	TokenBootstrap    Secreto   `json:"bootstrap_token"`
	ExpiraBootstrapEn time.Time `json:"bootstrap_expires_at"`
}
