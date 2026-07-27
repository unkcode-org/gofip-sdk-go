package gofip

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type transporte struct {
	claveAPI     string
	urlBase      string
	ambienteARCA AmbienteARCA
	cliente      *http.Client
}

type cuerpoDatos[T any] struct {
	Datos T `json:"data"`
}

type Archivo struct {
	Contenido     []byte
	TipoContenido string
	Nombre        string
}

type OpcionesGuardado struct {
	Sobrescribir bool
}

func (t *transporte) json(
	ctx context.Context,
	metodo string,
	ruta string,
	cuerpo any,
	cabeceras map[string]string,
	incluirAmbiente bool,
	destino any,
) error {
	var lector io.Reader
	if cuerpo != nil {
		contenido, err := json.Marshal(cuerpo)
		if err != nil {
			return fmt.Errorf("serializar solicitud: %w", err)
		}
		lector = bytes.NewReader(contenido)
	}

	req, err := http.NewRequestWithContext(ctx, metodo, t.urlBase+ruta, lector)
	if err != nil {
		return fmt.Errorf("crear solicitud: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.claveAPI)
	if cuerpo != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if incluirAmbiente {
		req.Header.Set("X-ARCA-Ambiente", string(t.ambienteARCA))
	}
	for nombre, valor := range cabeceras {
		if valor != "" {
			req.Header.Set(nombre, valor)
		}
	}

	resp, err := t.cliente.Do(req)
	if err != nil {
		return fmt.Errorf("realizar solicitud: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return errorDesdeRespuesta(resp)
	}
	if destino == nil || resp.StatusCode == http.StatusNoContent {
		_, _ = io.Copy(io.Discard, resp.Body)
		return nil
	}
	if err := json.NewDecoder(resp.Body).Decode(destino); err != nil {
		return fmt.Errorf("decodificar respuesta JSON: %w", err)
	}
	return nil
}

func (t *transporte) texto(ctx context.Context, ruta string, autenticar bool) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, t.urlBase+ruta, nil)
	if err != nil {
		return "", fmt.Errorf("crear solicitud: %w", err)
	}
	req.Header.Set("Accept", "text/plain, application/json")
	if autenticar {
		req.Header.Set("Authorization", "Bearer "+t.claveAPI)
	}

	resp, err := t.cliente.Do(req)
	if err != nil {
		return "", fmt.Errorf("realizar solicitud: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", errorDesdeRespuesta(resp)
	}
	contenido, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("leer respuesta: %w", err)
	}

	if strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		var valor string
		if err := json.Unmarshal(contenido, &valor); err != nil {
			return "", fmt.Errorf("decodificar respuesta JSON: %w", err)
		}
		return valor, nil
	}
	return string(contenido), nil
}

func (t *transporte) binario(
	ctx context.Context,
	metodo string,
	ruta string,
	cuerpo any,
	cabeceras map[string]string,
) (*Archivo, error) {
	var lector io.Reader
	if cuerpo != nil {
		contenido, err := json.Marshal(cuerpo)
		if err != nil {
			return nil, fmt.Errorf("serializar solicitud: %w", err)
		}
		lector = bytes.NewReader(contenido)
	}

	req, err := http.NewRequestWithContext(ctx, metodo, t.urlBase+ruta, lector)
	if err != nil {
		return nil, fmt.Errorf("crear solicitud: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+t.claveAPI)
	req.Header.Set("Accept", "application/pdf, text/html")
	if cuerpo != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for nombre, valor := range cabeceras {
		if valor != "" {
			req.Header.Set(nombre, valor)
		}
	}

	resp, err := t.cliente.Do(req)
	if err != nil {
		return nil, fmt.Errorf("realizar solicitud: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errorDesdeRespuesta(resp)
	}
	return archivoDesdeRespuesta(resp)
}

func (t *transporte) descargar(ctx context.Context, urlFirmada string) (*Archivo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlFirmada, nil)
	if err != nil {
		return nil, fmt.Errorf("crear descarga: %w", err)
	}
	req.Header.Set("Accept", "application/pdf, text/html, application/octet-stream")

	resp, err := t.cliente.Do(req)
	if err != nil {
		return nil, fmt.Errorf("descargar archivo: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, errorDesdeRespuesta(resp)
	}
	return archivoDesdeRespuesta(resp)
}

func archivoDesdeRespuesta(resp *http.Response) (*Archivo, error) {
	contenido, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("leer archivo: %w", err)
	}

	nombre := ""
	if disposicion := resp.Header.Get("Content-Disposition"); disposicion != "" {
		if _, parametros, err := mime.ParseMediaType(disposicion); err == nil {
			nombre = filepath.Base(parametros["filename"])
		}
	}
	if nombre == "" && resp.Request != nil && resp.Request.URL != nil {
		candidato := filepath.Base(resp.Request.URL.Path)
		if candidato != "." && candidato != string(filepath.Separator) {
			nombre = candidato
		}
	}

	return &Archivo{
		Contenido:     contenido,
		TipoContenido: resp.Header.Get("Content-Type"),
		Nombre:        nombre,
	}, nil
}

func errorDesdeRespuesta(resp *http.Response) error {
	contenido, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("leer respuesta de error HTTP %d: %w", resp.StatusCode, err)
	}

	var detalle ErrorAPI
	if err := json.Unmarshal(contenido, &detalle); err == nil &&
		(detalle.Titulo != "" || detalle.Detalle != "" || detalle.Codigo != "") {
		if detalle.Estado == 0 {
			detalle.Estado = resp.StatusCode
		}
		return &detalle
	}

	cuerpo := strings.TrimSpace(string(contenido))
	if cuerpo == "" {
		cuerpo = http.StatusText(resp.StatusCode)
	}
	return &ErrorHTTP{Estado: resp.StatusCode, Cuerpo: cuerpo}
}

func guardarArchivo(archivo *Archivo, ruta string, opciones OpcionesGuardado) (retErr error) {
	if archivo == nil {
		return errors.New("el archivo es obligatorio")
	}
	if strings.TrimSpace(ruta) == "" {
		return errors.New("la ruta destino es obligatoria")
	}

	ruta = filepath.Clean(ruta)
	if err := os.MkdirAll(filepath.Dir(ruta), 0o755); err != nil {
		return fmt.Errorf("crear directorio destino: %w", err)
	}

	banderas := os.O_WRONLY | os.O_CREATE
	if opciones.Sobrescribir {
		banderas |= os.O_TRUNC
	} else {
		banderas |= os.O_EXCL
	}
	f, err := os.OpenFile(ruta, banderas, 0o644)
	if err != nil {
		return fmt.Errorf("crear archivo destino: %w", err)
	}

	completo := false
	defer func() {
		if err := f.Close(); err != nil {
			if retErr == nil {
				retErr = fmt.Errorf("cerrar archivo destino: %w", err)
			}
			completo = false
		}
		if !completo {
			_ = os.Remove(ruta)
		}
	}()

	if _, err := f.Write(archivo.Contenido); err != nil {
		return fmt.Errorf("escribir archivo destino: %w", err)
	}
	completo = true
	return nil
}
