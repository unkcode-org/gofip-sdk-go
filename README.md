# GOFIP SDK para Go

Cliente idiomático en Go para consumir la API v2 de GOFIP y operar con los
servicios de ARCA. El SDK centraliza la autenticación con API key, la selección
del ambiente de ARCA, la serialización de solicitudes y respuestas, el manejo
de errores HTTP y la descarga de documentos.

## Versiones y contrato

| Componente | Versión usada |
| --- | --- |
| Módulo Go | `v2` (`github.com/unkcode-org/gofip-sdk-go/v2`) |
| API GOFIP | `2.0` y rutas bajo `/api/v2` |
| Especificación OpenAPI | `3.0.3` |
| Go | `1.22.2` o posterior |

El contrato HTTP de referencia está en
[`docs/openapi.yaml`](docs/openapi.yaml). Los campos `openapi: 3.0.3` e
`info.version: "2.0"` indican, respectivamente, la versión del formato OpenAPI
y la versión del contrato de la API. No deben confundirse con una versión
publicada del módulo Go.

## Funcionalidades

- Estado y disponibilidad de los servicios.
- Catálogos de ARCA: puntos de venta, tributos, conceptos, alícuotas, tipos de
  comprobante, opcionales, tipos de documento, monedas y cotizaciones.
- Consultas de padrón P4, P5, P10 y P13.
- Emisión de facturas y consulta de comprobantes.
- Generación y descarga de comprobantes, facturas y recibos en PDF.
- Errores tipados para respuestas de la API y respuestas HTTP no estructuradas.
- Cliente HTTP y URL base configurables.

El archivo OpenAPI describe también operaciones administrativas y de
autenticación del panel web que no necesariamente están expuestas por este
SDK. La superficie pública implementada se encuentra en los servicios
`Utilidades`, `Catalogos`, `Padron`, `Facturacion` y `Documentos`.

## Instalación

```bash
go get github.com/unkcode-org/gofip-sdk-go/v2@latest
```

## Configuración

La API key es obligatoria. El cliente la envía como token Bearer y, para las
operaciones de ARCA, agrega el ambiente seleccionado en la cabecera
`X-ARCA-Ambiente`.

```go
cliente, err := gofip.NuevoCliente(gofip.OpcionesCliente{
    ClaveAPI:     os.Getenv("GOFIP_API_KEY"),
    AmbienteARCA: gofip.AmbienteProduccion,
})
if err != nil {
    log.Fatal(err)
}
```

Opciones disponibles:

| Opción | Descripción | Valor predeterminado |
| --- | --- | --- |
| `ClaveAPI` | API key de GOFIP. Es obligatoria. | — |
| `URLBase` | URL absoluta del servicio. | `https://arca.unkcode.com` |
| `AmbienteARCA` | `AmbienteProduccion` o `AmbienteHomologacion`. | `AmbienteProduccion` |
| `ClienteHTTP` | `*http.Client` personalizado. | Timeout de 90 segundos |

Para evitar exponer credenciales, usá variables de entorno y no incluyas la
API key en el código fuente:

```bash
export GOFIP_API_KEY="tu-api-key"
```

## Uso rápido

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    gofip "github.com/unkcode-org/gofip-sdk-go/v2"
)

func main() {
    cliente, err := gofip.NuevoCliente(gofip.OpcionesCliente{
        ClaveAPI: os.Getenv("GOFIP_API_KEY"),
    })
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancelar := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancelar()

    monedas, err := cliente.Catalogos.ListarMonedas(ctx)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%+v\n", monedas.Monedas)
}
```

## Servicios disponibles

| Servicio | Operaciones principales |
| --- | --- |
| `cliente.Utilidades` | Salud, disponibilidad y medios de pago |
| `cliente.Catalogos` | Estado de ARCA, catálogos y cotización de monedas |
| `cliente.Padron` | Consultas P4, P5, P10 y P13 |
| `cliente.Facturacion` | Emisión de facturas y consulta de comprobantes |
| `cliente.Documentos` | Generación, obtención, descarga y guardado de documentos |

Los tipos de cada dominio viven en paquetes separados:

```go
import (
    "github.com/unkcode-org/gofip-sdk-go/v2/modelos/catalogos"
    "github.com/unkcode-org/gofip-sdk-go/v2/modelos/documentos"
    "github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
    "github.com/unkcode-org/gofip-sdk-go/v2/modelos/padron"
)
```

## Emisión de una factura

```go
respuesta, err := cliente.Facturacion.EmitirFactura(
    ctx,
    facturacion.SolicitudEmitirFactura{
        Factura: facturacion.Comprobante{
            TipoComprobante: facturacion.ComprobanteFacturaC,
            Concepto:        facturacion.ConceptoProductos,
            TipoDocumento:   facturacion.DocumentoConsumidorFinal,
            CondicionFiscal: facturacion.CondicionConsumidorFinal,
            Detalles: []facturacion.ItemComprobante{{
                ID:             "SERV-001",
                Descripcion:    "Servicio profesional",
                PrecioUnitario: 12500.50,
                Cantidad:       1,
                AlicuotasIVA:   facturacion.AlicuotaCero,
                Subtotal:       12500.50,
            }},
        },
    },
    gofip.OpcionesEmisionFactura{
        ClaveIdempotencia: "factura-pedido-123",
        ReferenciaExterna: "pedido-123",
    },
)
if err != nil {
    log.Fatal(err)
}
```

Usá una `ClaveIdempotencia` estable y única por operación de negocio. Si una
solicitud debe reintentarse por un error de red, reutilizá la misma clave para
evitar emisiones duplicadas. `ReferenciaExterna` permite asociar la operación
con un identificador propio.

Hay ejemplos completos en [`examples/`](examples/):

- [`cliente`](examples/cliente/main.go)
- [`catalogos`](examples/catalogos/main.go)
- [`padron`](examples/padron/main.go)
- [`emitir-factura`](examples/emitir-factura/main.go)
- [`generar-comprobante`](examples/generar-comprobante/main.go)
- [`generar-recibo`](examples/generar-recibo/main.go)

Podés ejecutarlos desde la raíz del repositorio:

```bash
go run ./examples/catalogos
```

## Manejo de errores

Las respuestas de error estructuradas se devuelven como `*gofip.ErrorAPI`. Si
el servidor responde con un cuerpo que no sigue ese formato, el SDK devuelve
`*gofip.ErrorHTTP`.

```go
var errorAPI *gofip.ErrorAPI
var errorHTTP *gofip.ErrorHTTP

switch {
case errors.As(err, &errorAPI):
    log.Printf(
        "GOFIP respondió %d (%s): %s [request_id=%s]",
        errorAPI.Estado,
        errorAPI.Codigo,
        errorAPI.Detalle,
        errorAPI.IDSolicitud,
    )
case errors.As(err, &errorHTTP):
    log.Printf("respuesta HTTP %d: %s", errorHTTP.Estado, errorHTTP.Cuerpo)
case err != nil:
    log.Printf("error de transporte: %v", err)
}
```

Todos los métodos reciben un `context.Context`. En aplicaciones reales es
recomendable definir deadlines o cancelación y configurar un `http.Client`
acorde a las necesidades del servicio.

## Descarga y guardado de archivos

Los métodos de descarga devuelven un `*gofip.Archivo` con el contenido, el tipo
MIME y el nombre informado por el servidor. Las variantes
`DescargarYGuardar...` y `GenerarYGuardar...` crean los directorios necesarios.
Por seguridad, no sobrescriben un archivo existente salvo que se indique:

```go
err := cliente.Documentos.DescargarYGuardarPDFRecibo(
    ctx,
    respuesta.PDFURL,
    "documentos/recibo.pdf",
    gofip.OpcionesGuardado{Sobrescribir: true},
)
```

## Desarrollo

Cloná el repositorio y ejecutá:

```bash
go test ./...
go vet ./...
```

La organización del proyecto se detalla en
[`docs/folder-structure.md`](docs/folder-structure.md). Ante un cambio del
contrato, actualizá `docs/openapi.yaml` junto con los modelos, servicios,
pruebas y ejemplos afectados, manteniendo sincronizadas las versiones
indicadas al comienzo de este README.

## Licencia

Este proyecto se distribuye bajo la
[Apache License 2.0](LICENSE).
