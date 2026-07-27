package gofip_test

import (
	"encoding/json"
	"testing"

	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/catalogos"
	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
	"github.com/unkcode-org/gofip-sdk-go/v2/modelos/padron"
)

func TestUnmarshalCatalogoConIDAlicuotaString(t *testing.T) {
	var respuesta catalogos.RespuestaAlicuotas
	if err := json.Unmarshal(
		[]byte(`{"alicuotas":[{"id":"5","descripcion":"21%","fecha_desde":"20000101"}]}`),
		&respuesta,
	); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(respuesta.Alicuotas) != 1 ||
		respuesta.Alicuotas[0].ID == nil ||
		*respuesta.Alicuotas[0].ID != "5" {
		t.Fatalf("respuesta inesperada: %#v", respuesta)
	}
}

func TestUnmarshalPadronConTimestamp(t *testing.T) {
	var respuesta padron.RespuestaPersonaPadron
	if err := json.Unmarshal(
		[]byte(`{"persona_p13":{"id_persona":30123456789,"fecha_nacimiento":"1990-01-02T00:00:00Z"}}`),
		&respuesta,
	); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if respuesta.PersonaP13 == nil ||
		respuesta.PersonaP13.FechaNacimiento == nil ||
		respuesta.PersonaP13.FechaNacimiento.Year() != 1990 {
		t.Fatalf("respuesta inesperada: %#v", respuesta)
	}
}

func TestUnmarshalComprobanteConsultado(t *testing.T) {
	var comprobante facturacion.ComprobanteConsultado
	if err := json.Unmarshal(
		[]byte(`{
			"punto_venta":2,
			"tipo_comprobante":11,
			"comprobante_desde":9,
			"importe_total":12500.5,
			"alicuotas_iva":[{"id":3,"base_imponible":12500.5,"importe":0}]
		}`),
		&comprobante,
	); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if comprobante.PuntoVenta == nil ||
		*comprobante.PuntoVenta != 2 ||
		len(comprobante.AlicuotasIVA) != 1 {
		t.Fatalf("comprobante inesperado: %#v", comprobante)
	}
}
