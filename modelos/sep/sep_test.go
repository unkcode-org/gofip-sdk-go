package sep

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"testing"
)

func TestSecretoRedactaFormatosJSONYSlog(t *testing.T) {
	var s Secreto
	if err := json.Unmarshal([]byte(`"capacidad-privada"`), &s); err != nil {
		t.Fatal(err)
	}
	if s.Valor() != "capacidad-privada" {
		t.Fatal("acceso explícito incorrecto")
	}
	var out bytes.Buffer
	slog.New(slog.NewJSONHandler(&out, nil)).Info("diagnostico", "secret", s, "nested", struct{ Token Secreto }{s})
	data, _ := json.Marshal(struct{ Token Secreto }{s})
	texto := fmt.Sprintf("%v %+v %#v %s %q %s %s", s, s, s, s, s, data, out.String())
	if strings.Contains(texto, s.Valor()) {
		t.Fatal("secreto visible en diagnóstico")
	}
	if err := json.Unmarshal([]byte(`{}`), &s); err == nil {
		t.Fatal("secreto no string aceptado")
	}
}
