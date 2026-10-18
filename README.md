# GOFIP SDK para Go

Cliente idiomático en Go para consumir la API v2 de GOFIP y operar con los
servicios de ARCA. El SDK centraliza la autenticación con API key, la selección
del ambiente de ARCA, la serialización de solicitudes y respuestas, el manejo
de errores HTTP y la descarga de documentos.

## Versiones y contrato

| Componente | Versión usada |
| --- | --- |
| Módulo Go | `v2` (`github.com/unkcode-org/gofip-sdk-go/v2`) |
| API GOFIP | `3.0.0` y rutas bajo `/api/v2` |
| Especificación OpenAPI | `3.0.3` |
| Go | `1.22.2` o posterior |

El contrato HTTP de referencia está en
[`docs/openapi.yaml`](docs/openapi.yaml). Los campos `openapi: 3.0.3` e
`info.version: "3.0.0"` indican, respectivamente, la versión del formato OpenAPI
y la versión del contrato de la API. No deben confundirse con una versión
publicada del módulo Go.

## Funcionalidades

- Estado y disponibilidad de los servicios.
- Catálogos de ARCA: puntos de venta, tributos, conceptos, alícuotas, tipos de
  comprobante, opcionales, tipos de documento, monedas y cotizaciones.
- Consultas de padrón P4, P5, P10 y P13.
- Emisión de facturas y consulta de comprobantes.
- Importación fiscal histórica verificada de GOFIP Legacy, con API key normal del mismo emisor/PV.
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
    log.Fatal("no se pudo completar la operación GOFIP")
}
```

Opciones disponibles:

| Opción | Descripción | Valor predeterminado |
| --- | --- | --- |
| `ClaveAPI` | API key de GOFIP. Es obligatoria. | — |
| `URLBase` | URL absoluta del servicio. | `https://arca.apps.staging.unkcode.com` |
| `AmbienteARCA` | `AmbienteProduccion` o `AmbienteHomologacion`. | `AmbienteProduccion` |
| `ClienteHTTP` | `*http.Client` personalizado. | Timeout de 90 segundos |

### Cliente HTTP instrumentado

Podés pasar un `http.Client` propio mediante `ClienteHTTP`. El SDK lo utiliza
para todas las solicitudes, incluidas las descargas de documentos por URL
firmada, por lo que tus middlewares de logs, métricas u OpenTelemetry observan
el tráfico completo. Al proporcionar un cliente personalizado, configurá en él
el timeout, transporte y demás políticas: el timeout predeterminado de 90
segundos solo se aplica cuando `ClienteHTTP` es `nil`.

El siguiente ejemplo envuelve el transporte para registrar la duración de cada
solicitud. Podés reemplazarlo por cualquier transporte instrumentado, como el
que ofrece `otelhttp`:

```go
type transporteConMetricas struct {
    siguiente http.RoundTripper
}

func (t transporteConMetricas) RoundTrip(solicitud *http.Request) (*http.Response, error) {
    inicio := time.Now()
    transporte := t.siguiente
    if transporte == nil {
        transporte = http.DefaultTransport
    }

    respuesta, err := transporte.RoundTrip(solicitud)
    log.Printf("gofip: método=%s duración=%s error=%t", solicitud.Method, time.Since(inicio), err != nil)
    return respuesta, err
}

clienteHTTP := &http.Client{
    Timeout: 30 * time.Second,
    Transport: transporteConMetricas{
        siguiente: http.DefaultTransport,
    },
}

cliente, err := gofip.NuevoCliente(gofip.OpcionesCliente{
    ClaveAPI:    os.Getenv("GOFIP_API_KEY"),
    ClienteHTTP: clienteHTTP,
})
if err != nil {
    log.Fatal("no se pudo completar la operación GOFIP")
}
```

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
        log.Fatal("no se pudo completar la operación GOFIP")
    }

    ctx, cancelar := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancelar()

    monedas, err := cliente.Catalogos.ListarMonedas(ctx)
    if err != nil {
        log.Fatal("no se pudo completar la operación GOFIP")
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
    log.Fatal("no se pudo completar la operación GOFIP")
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
- [`generar-pdf-factura`](examples/generar-pdf-factura/main.go)
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
    log.Printf("respuesta HTTP %d", errorHTTP.Estado)
case err != nil:
    log.Printf("error de transporte")
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

## Facturas persistidas: CAE y documento ID

Estas operaciones usan la API key Bearer y leen exclusivamente MySQL. Una clave
puede consultar las facturas de **todos los puntos de venta de su mismo emisor**,
sin ampliar el alcance de emisión. El emisor se resuelve en GOFIP; el SDK no lo
recibe como parámetro. Ambas consultas devuelven `facturacion.RespuestaFacturaPersistida`.

```go
factura, err := cliente.Facturacion.ObtenerFacturaPorCAE(ctx, cae)
// O utilizar el identificador estable DocumentoID devuelto al emitir:
factura, err = cliente.Facturacion.ObtenerFacturaPorId(ctx, documentoID)
```

Un CAE inexistente o inaccesible produce `*gofip.ErrorAPI` con estado 404.
Un CAE ambiguo dentro del emisor produce 409, código `13006`. Se desambigua sin
seleccionar una factura arbitraria:

```go
puntoVenta := 4
tipo := facturacion.ComprobanteFacturaA
numero := int64(125)
factura, err := cliente.Facturacion.ObtenerFacturaPorCAEConOpciones(
    ctx, cae, gofip.OpcionesConsultaFacturaPorCAE{
        PuntoVenta: &puntoVenta, TipoComprobante: &tipo,
        NumeroComprobante: &numero,
        // DocumentoID también puede restringir la búsqueda.
    },
)
```

`ObtenerComprobante(ctx, puntoVenta, tipo, numero)` conserva la consulta directa
a ARCA. `ObtenerFacturaPorCAE` y `ObtenerFacturaPorId` nunca usan ARCA como fallback.
El ID es el del documento GOFIP, no el de la fila interna `factura_id`.

La respuesta contiene sólo datos persistidos. No reconstruye datos históricos del
receptor desde la ficha mutable ni recupera SKU originales. Las fechas fiscales
son strings; los timestamps son `time.Time`. `pdf.disponible` refleja metadata y
no comprueba la existencia del objeto remoto ni genera una URL firmada.

**Ambiente:** GOFIP no persiste el ambiente en Documento/Factura. Estas consultas
no envían `X-ARCA-Ambiente` y no filtran por homologación/producción. Los registros
exclusivos de Legacy se incorporan mediante `ImportarFactura`, previa verificación
obligatoria en ARCA Producción y API key normal del mismo emisor/PV. Ver
[importación histórica](docs/importacion-facturas.md) para modelos, uso, errores,
idempotencia y límites de PDF. No hay batch ni emisión como fallback.

## Integraciones SEP: un cliente y una credencial distintos

```go
import sepmodel "github.com/unkcode-org/gofip-sdk-go/v2/modelos/sep"

clienteSEP, err := gofip.NuevoClienteSEP(gofip.OpcionesClienteSEP{
    ClaveIntegracion: os.Getenv("GOFIP_SEP_KEY"),
    URLBase: os.Getenv("GOFIP_API_URL"), // obligatoria, sin valor por defecto
    ClienteHTTP: &http.Client{Timeout: 15 * time.Second},
})
if err != nil { return err }
solicitud := sepmodel.SolicitudCrearSesionEmbedded{
    Permisos: []sepmodel.PermisoEmbedded{sepmodel.PermisoFacturasLeer},
}
sesion, err := clienteSEP.CrearSesion(ctx, solicitud)
// Ante una recuperación, solicitar un bootstrap nuevo:
bootstrap, err := clienteSEP.CrearBootstrap(ctx, solicitud)
```

| Credencial | Cliente | Uso |
| --- | --- | --- |
| `GOFIP_API_KEY` | `NuevoCliente` | Bearer; emisión, ARCA, facturas persistidas y documentos |
| `GOFIP_SEP_KEY` | `NuevoClienteSEP` | IntegrationKey; iniciar y recuperar acceso Embedded |

SEP envía exclusivamente `Authorization: IntegrationKey ...`, sin ambiente ARCA.
Crea requests `data` y decodifica responses `data`. Los permisos se validan contra
los doce valores de OpenAPI: no admite listas vacías, valores desconocidos ni
duplicados. Sólo el backend decide si caben dentro del techo autorizado.
`invoices:create` es un valor del enum y no agrega un endpoint de emisión Embedded.

`CrearSesion` devuelve URL de documento, vencimientos y bootstrap. Todavía no es
una sesión de navegador. `CrearBootstrap` emite un bootstrap nuevo de recuperación;
no recupera un secreto anterior ni crea otro documento. El canje, iframe,
postMessage, Origin/CSRF y almacenamiento del navegador pertenecen al frontend.

Los tokens y `URLDocumento` usan `sepmodel.Secreto`. `fmt`, `slog` y JSON los
redactan. Para entregar una capacidad deliberadamente al frontend, usá
`sesion.Datos.TokenBootstrap.Valor()` y `sesion.Datos.URLDocumento.Valor()` en un DTO
propio, dentro de un canal autorizado y sin logs/cache. Serializar directamente la
respuesta del SDK **entrega valores redactados**. No registrar structs completos,
headers, cuerpos o URLs de diagnóstico desde instrumentación HTTP personalizada.
La IntegrationKey permanece exclusivamente en el backend del integrador.

SEP bloquea redirecciones, no usa cookies y no reintenta automáticamente los POST.
Aplica estas políticas sobre una copia de `http.Client`; conserva el transporte y
el timeout recibidos y no modifica el cliente original. Sin cliente personalizado,
el timeout es 90 segundos. Los errores remotos se sanitizan conservando los tipos
`ErrorAPI`/`ErrorHTTP`, el estado y códigos numéricos seguros. Los errores de un
transporte personalizado no se exponen crudos; cancelaciones y deadlines conservan
`errors.Is(err, context.Canceled/DeadlineExceeded)`.

## Límites de idempotencia

La semántica de emisión no cambia: usá la misma credencial, clave estable y factura
para reintentar la misma operación comercial. Sin clave explícita no hay deduplicación
comercial estable. La referencia externa no es única; cambiar de credencial cambia
el ámbito de deduplicación. El hash del backend cubre el comprobante, no ambiente,
metadatos, papel ni referencia externa. Un intento pendiente/error puede devolver
conflicto. No reenviar con una clave nueva tras un resultado fiscal incierto.
La consulta por CAE no reconcilia autorizaciones que no llegaron a persistirse.

## Verificación y publicación

```bash
go test ./...
go vet ./...
go test -race ./...
```

Ver [arquitectura](docs/arquitectura.md) y [validación manual](docs/validacion.md).
El contrato de referencia es OpenAPI 3.0.0 del backend. No se publica una versión
como parte de estas pruebas. La secuencia prevista es `v2.0.0-rc.1` → `v2.0.0`,
con imports `/v2`, validación externa y aprobación de publicación.

SEP deriva el emisor exclusivamente de la Integration Key. CrearSesion y CrearBootstrap reciben solo permisos; la integración pertenece a un único emisor.
