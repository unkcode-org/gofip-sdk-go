package gofip

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OpcionesClienteSEP struct {
	ClaveIntegracion string `json:"-"`
	URLBase          string
	ClienteHTTP      *http.Client `json:"-"`
}

func (OpcionesClienteSEP) MarshalJSON() ([]byte, error) {
	return []byte(`{"configuracion":"[REDACTADO]"}`), nil
}

func (OpcionesClienteSEP) Format(f fmt.State, _ rune) {
	_, _ = fmt.Fprint(f, "OpcionesClienteSEP{[REDACTADO]}")
}
func (OpcionesClienteSEP) LogValue() slog.Value {
	return slog.StringValue("OpcionesClienteSEP{[REDACTADO]}")
}

// ClienteSEP sólo autentica operaciones servidor a servidor con IntegrationKey.
// No crea una sesión de navegador ni implementa el canje de bootstrap.
type ClienteSEP struct {
	clave   string
	urlBase string
	cliente *http.Client
}

func (*ClienteSEP) Format(f fmt.State, _ rune) { _, _ = fmt.Fprint(f, "ClienteSEP{[REDACTADO]}") }
func (*ClienteSEP) LogValue() slog.Value       { return slog.StringValue("ClienteSEP{[REDACTADO]}") }

func NuevoClienteSEP(o OpcionesClienteSEP) (*ClienteSEP, error) {
	clave := strings.TrimSpace(o.ClaveIntegracion)
	if clave == "" || strings.ContainsAny(clave, " \t\r\n") {
		return nil, errors.New("la clave de integración es obligatoria y no debe contener espacios")
	}
	u, err := url.Parse(o.URLBase)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.User != nil || u.RawQuery != "" || u.ForceQuery || u.Fragment != "" || strings.Contains(o.URLBase, "#") {
		return nil, errors.New("SEP requiere URL base HTTP(S) absoluta explícita, sin credenciales, query ni fragmento")
	}
	var cliente http.Client
	if o.ClienteHTTP == nil {
		cliente = http.Client{Timeout: 90 * time.Second}
	} else {
		cliente = *o.ClienteHTTP
	}
	// Políticas locales sobre una copia: jamás modificar al consumidor ni usar
	// cookies para autenticar la superficie M2M.
	cliente.Jar = nil
	cliente.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return &ClienteSEP{clave: clave, urlBase: strings.TrimRight(o.URLBase, "/"), cliente: &cliente}, nil
}
