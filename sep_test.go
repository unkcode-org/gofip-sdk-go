package gofip

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

const claveSEPPrueba = "gofip_int_public_secreto-privado"
const tokenSEPPrueba = "gofip_emb_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
const urlDocumentoSEPPrueba = "https://embedded.example/embed/capacidad-privada"

func TestNuevoClienteSEPConfiguracion(t *testing.T) {
	for _, o := range []OpcionesClienteSEP{
		{}, {ClaveIntegracion: claveSEPPrueba}, {ClaveIntegracion: claveSEPPrueba, URLBase: "/relativa"},
		{ClaveIntegracion: claveSEPPrueba, URLBase: "ftp://example.com"},
		{ClaveIntegracion: claveSEPPrueba, URLBase: "https://user:secret@example.com"},
		{ClaveIntegracion: claveSEPPrueba, URLBase: "https://example.com?token=privado"},
		{ClaveIntegracion: claveSEPPrueba, URLBase: "https://example.com#privado"},
		{ClaveIntegracion: "clave\r\ninyectada", URLBase: "https://example.com"},
	} {
		if _, err := NuevoClienteSEP(o); err == nil {
			t.Fatal("configuración inválida aceptada")
		}
	}
	o := OpcionesClienteSEP{ClaveIntegracion: claveSEPPrueba, URLBase: "https://example.com/"}
	c, err := NuevoClienteSEP(o)
	if err != nil {
		t.Fatal(err)
	}
	if c.cliente.Timeout != 90*time.Second {
		t.Fatal("timeout predeterminado incorrecto")
	}
	if strings.Contains(fmt.Sprintf("%v %+v %#v", o, o, o), claveSEPPrueba) {
		t.Fatal("opciones exponen credencial")
	}
	if strings.Contains(fmt.Sprintf("%v %+v %#v", c, c, c), claveSEPPrueba) {
		t.Fatal("cliente expone credencial")
	}
	b, _ := json.Marshal(o)
	if strings.Contains(string(b), claveSEPPrueba) {
		t.Fatal("opciones JSON exponen credencial")
	}
}

func TestSEPCreacionYRecuperacionContrato(t *testing.T) {
	var llamadas atomic.Int32
	servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		llamadas.Add(1)
		if r.Method != "POST" || r.Header.Get("Authorization") != "IntegrationKey "+claveSEPPrueba {
			t.Error("autenticación incorrecta")
		}
		if r.Header.Get("Cookie") != "" || r.Header.Get("X-ARCA-Ambiente") != "" {
			t.Error("SEP mezcló autenticaciones o ambiente")
		}
		if r.Header.Get("Accept") != "application/json" || r.Header.Get("Content-Type") != "" {
			t.Error("cabeceras JSON")
		}
		b, err := io.ReadAll(r.Body)
		if err != nil || len(b) != 0 || r.URL.RawQuery != "" {
			t.Error("emisión SEP debe tener cuerpo y query vacíos")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(201)
		switch r.URL.Path {
		case "/api/v2/embedded/sessions":
			_, _ = fmt.Fprintf(w, `{"data":{"document_url":%q,"document_expires_at":"2026-10-16T12:30:00Z","bootstrap_token":%q,"bootstrap_expires_at":"2026-10-16T12:01:00Z"}}`, urlDocumentoSEPPrueba, tokenSEPPrueba)
		case "/api/v2/embedded/bootstraps":
			_, _ = fmt.Fprintf(w, `{"data":{"bootstrap_token":%q,"bootstrap_expires_at":"2026-10-16T12:01:00Z"}}`, tokenSEPPrueba)
		default:
			t.Error("ruta inesperada")
		}
	}))
	defer servidor.Close()
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(servidor.URL)
	jar.SetCookies(u, []*http.Cookie{{Name: "jwt", Value: "no-enviar"}})
	original := servidor.Client()
	original.Jar = jar
	original.Timeout = 7 * time.Second
	original.CheckRedirect = func(*http.Request, []*http.Request) error { return errors.New("politica-original") }
	c, err := NuevoClienteSEP(OpcionesClienteSEP{ClaveIntegracion: claveSEPPrueba, URLBase: servidor.URL + "/", ClienteHTTP: original})
	if err != nil {
		t.Fatal(err)
	}
	sesion, err := c.CrearSesion(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if sesion.Datos.TokenBootstrap.Valor() != tokenSEPPrueba || sesion.Datos.URLDocumento.Valor() != urlDocumentoSEPPrueba || sesion.Datos.ExpiraDocumentoEn.Year() != 2026 {
		t.Fatal("respuesta incorrecta")
	}
	bootstrap, err := c.CrearBootstrap(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if bootstrap.Datos.TokenBootstrap.Valor() != tokenSEPPrueba {
		t.Fatal("bootstrap incorrecto")
	}
	if llamadas.Load() != 2 {
		t.Fatal("reintentos inesperados")
	}
	if original.Jar != jar || original.Timeout != 7*time.Second || original.CheckRedirect(nil, nil).Error() != "politica-original" {
		t.Fatal("cliente original modificado")
	}
	for _, v := range []any{sesion, bootstrap, sesion.Datos.TokenBootstrap, sesion.Datos.URLDocumento} {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		diagnostico := fmt.Sprintf("%v %+v %#v %s", v, v, v, b)
		for _, secreto := range []string{tokenSEPPrueba, urlDocumentoSEPPrueba} {
			if strings.Contains(diagnostico, secreto) {
				t.Fatal("capacidad expuesta")
			}
		}
	}
}

func TestSEPErroresRedactadosSinReintentos(t *testing.T) {
	for _, status := range []int{400, 401, 403, 404, 500, 302} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			var llamadas atomic.Int32
			servidor := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				llamadas.Add(1)
				w.Header().Set("Location", "https://evil.example/"+tokenSEPPrueba)
				w.WriteHeader(status)
				_, _ = fmt.Fprintf(w, `{"code":"18017","title":%q,"detail":%q,"instance":%q,"request_id":%q}`, claveSEPPrueba, tokenSEPPrueba, urlDocumentoSEPPrueba, claveSEPPrueba)
			}))
			defer servidor.Close()
			c, _ := NuevoClienteSEP(OpcionesClienteSEP{ClaveIntegracion: claveSEPPrueba, URLBase: servidor.URL, ClienteHTTP: servidor.Client()})
			_, err := c.CrearSesion(context.Background())
			if err == nil {
				t.Fatal("falta error")
			}
			var api *ErrorAPI
			if !errors.As(err, &api) || api.Estado != status || api.Codigo != "18017" {
				t.Fatalf("tipo incorrecto: %T", err)
			}
			b, _ := json.Marshal(api)
			for _, secret := range []string{claveSEPPrueba, tokenSEPPrueba, urlDocumentoSEPPrueba} {
				if strings.Contains(fmt.Sprintf("%v %#v %s", err, api, b), secret) {
					t.Fatal("error expone secreto")
				}
			}
			if llamadas.Load() != 1 {
				t.Fatal("se reintentó el POST")
			}
		})
	}
}

func TestSEPCancelacionDeadlineYTransporte(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { _, _ = io.Copy(io.Discard, r.Body); <-r.Context().Done() }))
	defer s.Close()
	c, _ := NuevoClienteSEP(OpcionesClienteSEP{ClaveIntegracion: claveSEPPrueba, URLBase: s.URL, ClienteHTTP: s.Client()})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := c.CrearSesion(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelación %v", err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Millisecond)
	defer cancel()
	if _, err := c.CrearBootstrap(ctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("deadline %v", err)
	}
	var llamadas atomic.Int32
	c, _ = NuevoClienteSEP(OpcionesClienteSEP{ClaveIntegracion: claveSEPPrueba, URLBase: "https://example.test", ClienteHTTP: &http.Client{Transport: roundTripperPrueba(func(*http.Request) (*http.Response, error) {
		llamadas.Add(1)
		return nil, fmt.Errorf("fallo con %s", claveSEPPrueba)
	})}})
	_, err := c.CrearBootstrap(context.Background())
	if err == nil || strings.Contains(err.Error(), claveSEPPrueba) || llamadas.Load() != 1 {
		t.Fatal("error transporte inseguro o reintentos")
	}
}

func TestSEPRespuestaInvalidaYErrorHTTP(t *testing.T) {
	for _, body := range []string{"no-json", `{}`, `{"data":{"bootstrap_token":"secreto-malformado"}}`} {
		c, _ := NuevoClienteSEP(OpcionesClienteSEP{ClaveIntegracion: claveSEPPrueba, URLBase: "https://example.test", ClienteHTTP: &http.Client{Transport: roundTripperPrueba(func(r *http.Request) (*http.Response, error) {
			return &http.Response{StatusCode: 201, Body: io.NopCloser(strings.NewReader(body)), Header: make(http.Header), Request: r}, nil
		})}})
		if _, err := c.CrearBootstrap(context.Background()); err == nil || strings.Contains(err.Error(), "secreto-malformado") {
			t.Fatal("respuesta inválida o error inseguro")
		}
	}
	c, _ := NuevoClienteSEP(OpcionesClienteSEP{ClaveIntegracion: claveSEPPrueba, URLBase: "https://example.test", ClienteHTTP: &http.Client{Transport: roundTripperPrueba(func(r *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: 502, Body: io.NopCloser(strings.NewReader(claveSEPPrueba + tokenSEPPrueba)), Header: make(http.Header)}, nil
	})}})
	_, err := c.CrearBootstrap(context.Background())
	var httpErr *ErrorHTTP
	if !errors.As(err, &httpErr) || httpErr.Estado != 502 || strings.Contains(httpErr.Cuerpo, claveSEPPrueba) {
		t.Fatal("error HTTP inseguro")
	}
}

func TestSEPDeadlineDuranteLecturaDelCuerpo(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		w.WriteHeader(201)
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer s.Close()
	c, _ := NuevoClienteSEP(OpcionesClienteSEP{ClaveIntegracion: claveSEPPrueba, URLBase: s.URL, ClienteHTTP: s.Client()})
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	_, err := c.CrearBootstrap(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("deadline leyendo respuesta: %v", err)
	}
}
