# Validación de la integración

## Antes de ejecutar ejemplos

Definir `GOFIP_API_URL`, `GOFIP_API_KEY` y, sólo para SEP, `GOFIP_SEP_KEY` mediante
variables del entorno. No pegarlas en logs, comandos compartidos ni archivos
versionados. La URL corresponde al backend V2 actualizado con OpenAPI 2.18.0.

Los ejemplos de consulta no generan efectos fiscales. El ejemplo de emisión sí
emite: usar homologación y una venta de prueba deliberada.

## Consultas

- `GOFIP_CAE`: CAE persistido conocido. `go run ./examples/consultar-factura-cae`.
- `GOFIP_DOCUMENTO_ID`: ID del documento. `go run ./examples/consultar-factura-id`.
- Verificar igualdad de datos usando los dos mecanismos.
- Usar las claves de dos puntos de venta del mismo emisor: ambas deben leerlo.
- Usar una clave de otro emisor: debe responder 404.
- Un CAE duplicado dentro del emisor debe responder 409/código 13006. Usar
  `ObtenerFacturaPorCAEConOpciones` para desambiguar.
- La consulta ARCA sigue siendo `ObtenerComprobante`; estos GET nunca la invocan.

## SEP

Definir `GOFIP_SEP_EMISOR_EXTERNO` con el mapping autorizado.

```bash
go run ./examples/sep-sesion
go run ./examples/sep-bootstrap
```

Los ejemplos sólo imprimen vencimientos. Crean capacidades reales y no realizan
el canje: no repetirlos automáticamente ni usar sus resultados de diagnóstico
como secretos. La integración del frontend debe recibir las capacidades por una
entrega explícita, autorizada y no cacheable.

## Suite y consumo externo

```bash
go test ./...
go vet ./...
go test -race ./...
```

Para probar desde otro módulo antes de publicar, agregar un `require` del módulo
`github.com/unkcode-org/gofip-sdk-go/v2` y un `replace` hacia este checkout; compilar
un consumidor que invoque las cuatro operaciones nuevas y los métodos existentes.
No cambiar el import a `/v0` ni publicar tags durante la validación.

El backend agrega una prueba MySQL opt-in:

```bash
go test ./internal/adapter/db/mysql -run '^TestConsultaFacturasMySQL$' -v
```

Requiere `GOFIP_FACTURAS_TEST_MYSQL_DSN` de un servidor de pruebas con permisos para
crear/eliminar una base temporal exclusiva. No usa tablas de la base del DSN.
La prueba real de GOFIP y MySQL queda a cargo del integrador antes de publicar.

El módulo externo de prueba está incluido en `integration/consumer`: ejecutar `go test ./...` desde esa carpeta. El replace local permite validar los imports `/v2` sin publicar tags.
