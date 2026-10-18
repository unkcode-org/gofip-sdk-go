# Importar autorizaciones históricas Legacy

`ImportarFactura` requiere Backend V2 con OpenAPI 2.19.0 y una API key normal
activa del Punto de Venta. La credencial sólo puede importar comprobantes de su
propio emisor y PV, incluso si ese emisor tiene varios PV. El consumidor no envía
emisor, ambiente ni autoridad adicional. El backend verifica la relación
`CredencialPuntoVenta -> PuntoVenta -> Emisor` y compara el número fiscal del PV
con `solicitud.Factura.PuntoVenta`. Otro PV devuelve 404 sin consultar ARCA.

```go
solicitud := facturacion.SolicitudImportarFactura{
    Factura: facturacion.FacturaHistorica{
        PuntoVenta: 3,
        TipoComprobante: facturacion.TipoComprobante(6),
        NumeroComprobante: 12345,
        CAE: "00123456789012",
    },
    Procedencia: facturacion.ProcedenciaFactura{Sistema: "gofip_legacy"},
}
resultado, err := cliente.Facturacion.ImportarFactura(ctx, solicitud,
    gofip.OpcionesImportacionFactura{
        ClaveIdempotencia: "historico:12345:original",
        ReferenciaExterna: "cliente:venta:123:factura:original",
    })
// Persistir resultado.DocumentoID en el cliente sólo si err == nil.
```

El backend consulta exclusivamente `FECompConsultar` en ARCA Producción,
compara el CAE esperado y persiste la autorización existente. Legacy siempre
pertenece a producción, incluso si `OpcionesCliente.AmbienteARCA` se configuró
para otras operaciones en homologación. El método no envía `X-ARCA-Ambiente` y
`OpcionesImportacionFactura` no ofrece selector. No crea una autorización, PDF,
correo ni fallback.

El SDK agrega el envelope `data`; las respuestas no usan envelope. 201 significa
`Importado=true`; 200 significa `YaExistia=true`. Ambos contienen el documento
persistido y evidencia de verificación. El origen de un documento ya existente
se conserva. Luego se usan `ObtenerFacturaPorId` y `ObtenerFacturaPorCAE`, con
lectura por todos los PV del mismo emisor y aislamiento entre emisores.

## Datos y precisión

`DatosHistoricos` es un conjunto opcional de afirmaciones fiscales; cada campo
enviado debe coincidir con la consulta. ARCA es la fuente fiscal: un dato ausente
no se completa con reglas actuales ni con valores por defecto. `RegistroID` en
procedencia es opcional. No exigir un ID central de factura que Legacy no posee.

Los nuevos decimales son `facturacion.DecimalHistorico`, strings exactos con
hasta 16 enteros y seis decimales. Por ejemplo, `"1234.123456"`; JSON numérico no
es válido para esos campos. Los punteros conservan ausencia. Para listas fiscales,
`nil` significa desconocida/ausente y un puntero a slice vacío significa lista
conocida vacía. El campo `Factura.Historico` conserva la respuesta fiscal exacta;
los campos float de lectura anteriores siguen siendo una proyección compatible.

`Detalle` conserva ítems y datos documentales históricos aportados por el
consumidor. Es información comercial que ARCA no certifica. No recalcula totales.
V1 conserva esa información y comunica completitud, pero devuelve un error
explícito al regenerar PDFs históricos con el renderer actual: 13106 si falta
información, 13107 si requiere el renderer histórico aún no implementado.

## Idempotencia, referencias y fallas

La clave operacional es obligatoria y estable (hasta 120 bytes). Su scope es
`emisor_id + credencial_id + punto_venta_id + clave`, separado de la emisión. La misma clave
con contenido distinto produce 409. Tras un timeout se puede repetir exactamente
la misma solicitud y clave; el SDK no reintenta automáticamente. Tampoco sigue
redirecciones, conserva cookies ni permite el replay transparente del body por
el transporte estándar. Un RoundTripper propio sigue siendo responsabilidad del
integrador y debe evitar retries y registros de secretos.

La identidad definitiva es `emisor + PV + tipo + número`, protegida por UNIQUE
real compartida con emisión normal. Distintas claves para la misma identidad
compatible devuelven el mismo documento. La referencia tiene scope
`emisor_id + credencial_id + punto_venta_id`; repetida hacia otro documento produce 409. No cambia la
unicidad de referencias de emisión. Ambos scopes usan SHA-256 de un array JSON
ordenado con los cuatro strings; los IDs aíslan credenciales, emisores y PV.
La referencia es opcional y no define la identidad fiscal. Una nueva credencial
crea otro scope, pero la identidad fiscal conserva el documento compatible.

| HTTP / código | Acción del consumidor |
|---|---|
| 400 / 13100 | Corregir solicitud inválida |
| 401 / 12004, 12005 | Revisar validez y expiración de la API key |
| 403 / 12003 | Credencial revocada o inactiva |
| 404 | Verificar que el PV sea el propio de la credencial |
| 404 / 13103 | ARCA no encontró el comprobante; investigar |
| 409 / 13101 | Resolver conflicto de identidad, datos, referencia o clave |
| 413 / 13100 | Reducir body a 1 MiB |
| 502 / 13105 | Revisar respuesta fiscal inválida/incompleta |
| 503 / 13104 | Reanudar con solicitud y clave estables al recuperarse ARCA |

Todos los errores conservan el tratamiento tipado `ErrorAPI`/`ErrorHTTP` para
respuestas HTTP, saneando diagnóstico remoto y cadenas de transporte. Los errores
locales son seguros. Cancelaciones y deadlines mantienen `errors.Is`. La respuesta
tiene límite de 2 MiB y el cliente HTTP es configurable; se copia su configuración
sin mutar la instancia original. Nunca convertir una falla de importación en
`EmitirFactura`.

## Ejemplo y validación

`examples/importar-factura` procesa una sola solicitud tipada sin envelope desde
`GOFIP_IMPORT_REQUEST_FILE`. Exige `GOFIP_IMPORT_IDEMPOTENCY_KEY`, `GOFIP_API_URL`
y `GOFIP_API_KEY`; `GOFIP_IMPORT_EXTERNAL_REFERENCE` es opcional. Ejecutarlo sólo
con datos y entorno autorizados. No se ejecuta como parte de tests.

```sh
go test ./...
go vet ./...
go test -race ./...
cd integration/consumer
go test ./...
```

Los tests usan servidores controlados y no importan datos ni autorizan facturas
reales. El backend aporta la prueba multi-proceso contra MariaDB/MySQL real y
un puerto que sólo expone consulta; este módulo no implementa batch ni contiene
adaptadores de Elixir u otros clientes.
