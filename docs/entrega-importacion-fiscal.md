# Entrega de importación fiscal histórica — actualizada 2026-10-18

Revisión de autorización del 2026-10-18: backend
`7918e9d373a8cbc298258cb9314d5d33e54d38b3`. La API key normal del propio emisor/PV
es toda la autorización necesaria. Esta revisión del SDK actualiza sólo textos,
comentarios, contrato copiado y su test; no cambia el método ni las opciones públicas.
La validación actual vuelve a pasar test/vet/race y el consumidor externo.
La evidencia completa de migración y concurrencia está en
`backend/docs/VALIDACION_IMPORTACION_FISCAL.md` del backend.

## 1. Commits y ubicación

Rama `codex/importacion-fiscal-historica` en ambos repositorios:

| Repositorio | Commit | Contenido |
|---|---|---|
| gofip-backend | `77a2292` | Especificación inicial en rojo |
| gofip-backend | `e9577b8790bc853f6ab486d01315ca9771b0fe0b` | Implementación, migración, tests y documentación |
| gofip-sdk-go/v2 | `fb0a31c` | Especificación inicial en rojo |
| gofip-sdk-go/v2 | Commit que incorpora este documento (`git log -1`) | Implementación, contrato sincronizado, ejemplo y validación |

Backend: `/Users/genarocivilotti/Code/gofip-backend/backend`.
SDK: `/Users/genarocivilotti/.codex/worktrees/importacion-fiscal-historica/sdk-go`.
El checkout original del SDK quedó intacto. No hubo push, PR, merge, tag ni release.

## 2. Archivos modificados

El inventario completo está al final. La implementación abarca dominio/puertos,
repositorios MySQL, adapter de consulta SOAP, handlers/ruta, errores, renderer
protegido y contrato/documentación. El SDK agrega método,
modelos tipados, ejemplo individual, tests y consumidor externo.

## 3. Migración

`202610160002_importacion_fiscal` registra las identidades de documentos
V2 autorizados en una transacción. Datos incompletos, PV inconsistentes o
identidades incompatibles hacen fallar el backfill con un error que identifica
el documento. Las constraints UNIQUE son la garantía final; no corrige datos.
Crea tres tablas: identidades, operaciones y referencias de importación.
Los campos de snapshot/verificación/completitud y la ampliación decimal se
incorporan al esquema histórico mediante el AutoMigrate existente. El rollback
rechaza datos fiscales/operacionales y conserva las tablas aditivas vacías.
No se agregan columnas ni backfills de ambiente.

## 4. Constraint fiscal

```sql
UNIQUE (emisor_id, punto_venta, tipo_comprobante, numero_comprobante)
UNIQUE (documento_id)
```

Son índices SQL reales en `identidades_fiscales`, compartidos por emisión e
importación. CAE se valida pero no reemplaza esa identidad.

## 5. Escritor normal V2

`CrearDocumentoFacturaContext` y `RegistrarFacturaAutorizadaContext` registran la
identidad en la misma transacción que la factura autorizada. Un documento pendiente
no reserva identidad. El contrato, hashes, retries y semántica de `EmitirFactura`
permanecen vigentes. Una colisión impide persistir una segunda identidad y hace
rollback, sin sobrescribir la autorización existente.

## 6. Endpoint

`POST /api/v2/facturas/importaciones`, Bearer API key, envelope `data`, request
estricto, 1 MiB máximo. Emisor resuelto desde autenticación; no acepta emisor ni
ambiente del body. Idempotency-Key obligatorio y referencia opcional.
201 crea; 200 devuelve existente; 409 para identidad/CAE/datos/referencia/clave
incompatibles; 404 ARCA ausente; 502 respuesta fiscal inválida; 503 indisponible.
Las facturas son documentos autorizados normales y usan los GET existentes por
ID/CAE, con lectura por emisor y aislamiento entre emisores.

## 7. Autorización

La API key normal activa del PV es suficiente. El backend recibe la credencial,
su PV y su emisor de la autenticación existente y exige que el PV solicitado sea
exactamente el asociado a esa credencial. Otro PV, aun del mismo emisor, devuelve
404 sin consultar ARCA. Se vuelve a validar la relación antes de persistir.
Claves y referencias se aíslan mediante SHA-256 del array JSON ordenado
`[emisor_id, credencial_id, punto_venta_id, clave_o_referencia]` en tablas separadas.
No hay configuración adicional de autorización para importaciones.

## 8. Consulta ARCA

Puerto de un solo método `ConsultarComprobanteProduccion`, independiente del
adapter amplio de emisión. Autenticación WSAA en contexto Producción y SOAP
`FECompConsultar` contra el endpoint oficial productivo. No expone autorización,
CAE/CAEA mutante, numeración ni selector de ambiente. `X-ARCA-Ambiente` no modifica
la política Legacy. No hay modo unverified ni fallback.
El CAE esperado y todas las afirmaciones fiscales enviadas se contrastan.
Los campos y estructuras recuperables conservan presencia, códigos y precisión.

## 9. API y modelos SDK

```go
cliente.Facturacion.ImportarFactura(ctx, solicitud, gofip.OpcionesImportacionFactura{
    ClaveIdempotencia: claveEstable,
    ReferenciaExterna: referenciaEstable,
})
```

Modelos principales en `modelos/facturacion/importacion.go`:
`SolicitudImportarFactura`, `FacturaHistorica`, `ProcedenciaFactura`,
`DatosFiscalesHistoricos`, `DecimalHistorico`, `DetalleHistoricoFactura`,
`ItemHistorico`, `DatosDocumentalesHistoricos`, `VerificacionImportacion`,
`CompletitudDocumental` y `RespuestaImportarFactura`, además de IVA/tributos/
asociados/opcionales/compradores/actividades/períodos/observaciones tipados.
`RespuestaFacturaPersistida` agrega verificación/completitud y snapshot histórico.
No hay maps públicos ni batch. Se preservan context, cliente configurable,
cancelación/deadline y errores HTTP tipados con redacción segura. No hay retries
automáticos ni redirects, cookies o replay de body por transporte estándar.

## 10. Tests ejecutados

Backend y SDK: `go test ./...`, `go vet ./...`, `go test -race ./...` y
`git diff --check`: PASS. Consumidor externo: test, vet y race: PASS.
SDK y consumidor externo con Go **1.22.2**:
`GOTOOLCHAIN=go1.22.2 go test -ldflags=-linkmode=external ./...`: PASS.
La ejecución directa del enlazador interno antiguo abortó en dyld por LC_UUID en
este macOS; el enlace externo permitió validar ejecución real, sin cambiar el
mínimo Go del módulo. No se oculta ese límite de la herramienta local.
La copia OpenAPI comprueba SHA-256/procedencia y ausencia de selector de autoridad.

## 11. Concurrencia real

MariaDB **13.0.2**, dos procesos, cinco goroutines por proceso, diez claves
operacionales distintas, mismo comprobante: **un documento, una factura, una
identidad, una referencia y el mismo documento_id en los diez resultados**.
Se ejecutó además backfill repetible y conservación SQL exacta de
`1234567890123456.123456`. Se verificaron índices reales en information_schema.
No se usó SQLite para probar la garantía multi-proceso.

## 12. Cero emisión

La prueba de arquitectura exige un puerto con exactamente un método de consulta.
El test de integración utiliza el adapter SOAP real con transporte controlado:
10 FECompConsultar y **0 FECAESolicitar / otras operaciones mutantes**.
Cualquier otra acción hace fallar la prueba. Autenticación falsa; sin secretos ni
llamadas reales a ARCA. La instancia de MariaDB fue descartable, sin red, y quedó
apagada. Se instaló MariaDB mediante Homebrew para ejecutar esta prueba; no se
activó su servicio del sistema.

## 13. OpenAPI

Contrato canónico **2.19.0**, endpoint, modelos, optional/nullable, autorización por credencial/PV,
procedencia, idempotencia, referencias, errores, origen y datos de consulta.
Copia SDK exacta con SHA-256
`82ed564929b6c46e608d7922dc0a4fb93e1f4e1c2d8ec7bc06d434b4975d92b1`.
Procedencia completa en `docs/CONTRATO.md`; test local detecta drift.

## 14. Documentación y observabilidad

Backend: `docs/IMPORTACION_FACTURAS_LEGACY.md` y
`docs/VALIDACION_IMPORTACION_FISCAL.md`. SDK: `docs/importacion-facturas.md`,
README/arquitectura/validación/procedencia y `examples/importar-factura`.
Eventos `invoice.import.started` / `invoice.import`, span `invoice.import` y
contador `gofip.invoices.imports`, con resultados acotados e identificadores
seguros. No incrementa métricas de emisión ni registra bodies, claves, tokens,
certificados o URLs firmadas.

## 15. Límites y riesgos pendientes

- No se contactó ARCA real ni se modificaron datos de producción. Certificados,
  delegaciones WSFE productivas y disponibilidad requieren validación posterior.
- Si el backfill encuentra datos incompatibles, aborta sin corregirlos.
  Drenar todos los escritores anteriores y detener escrituras fiscales durante
  el cambio; un binario antiguo no conoce la tabla de identidad compartida.
- PDF histórico queda explícitamente bloqueado en V1 (13106/13107). Se conservan
  datos aportados, sin reconstruir desde fichas actuales ni enviar correo.
- La verdad fiscal exacta está en `historico` y columnas DECIMAL ampliadas;
  los floats del modelo público anterior siguen siendo proyecciones compatibles.
- La verificación no recupera ítems, nombres ni domicilios históricos inexistentes.
  ARCA no certifica esos datos comerciales. Una diferencia fiscal produce conflicto.
- Las reservas fallidas se conservan para reanudar; definir retención operativa
  posteriormente. No se borran automáticamente ni se convierten en facturas.
- Una caída de DB tras una emisión normal sigue siendo un problema de conciliación
  preexistente: esta tarea no introduce retries ni autorización como recuperación.

## 16. Primer migrador real de Joyería Elixir, posteriormente

1. Revisar estos commits y coordinar despliegue/migración de esquema por separado.
2. Preparar configuración/delegación productiva del emisor.
3. Usar la API key normal del PV correspondiente a cada comprobante.
4. Implementar en Elixir el adaptador de `OrderPaymentFiscalData`, sin inventar
   ausencias ni recalcular importes; incluir los datos conocidos como afirmaciones.
5. Importar individualmente con clave/referencia estable, persistir documento_id
   y luego marcar la operación fiscal `issued`; guardar progreso y errores.
6. Legacy sin mapping significa migración pendiente y jamás habilita emisión.
   Mapping existente significa consultar V2. Recién después revisar la eliminación
   de `GOFIP_V2_ORDERS_FROM`.

No se modificaron Elixir, otros consumidores ni APIs Legacy. No hubo importaciones
reales, emisión, despliegue, migración de producción o publicación automática.

## Inventario Backend (rutas relativas a gofip-backend)

- `backend/README.md`
- `backend/cmd/main.go`
- `backend/docs/CHANGELOG.md`
- `backend/docs/IMPORTACION_FACTURAS_LEGACY.md`
- `backend/docs/VALIDACION_IMPORTACION_FISCAL.md`
- `backend/docs/importaciones_test.go`
- `backend/docs/openapi.yaml`
- `backend/internal/adapter/db/mysql/importacion_facturas_test.go`
- `backend/internal/adapter/db/mysql/importacion_http_test.go`
- `backend/internal/adapter/db/mysql/importacion_mysql_test.go`
- `backend/internal/adapter/db/mysql/integridad_fiscal.go`
- `backend/internal/adapter/db/mysql/integridad_fiscal_test.go`
- `backend/internal/adapter/db/mysql/migraciones.go`
- `backend/internal/adapter/db/mysql/migraciones_test.go`
- `backend/internal/adapter/db/mysql/repositorio_facturacion.go`
- `backend/internal/adapter/db/mysql/repositorio_importaciones.go`
- `backend/internal/adapter/db/mysql/repositorio_integracion_embedded_test.go`
- `backend/internal/adapter/db/mysql/repositorios.go`
- `backend/internal/adapter/db/mysql/repositorios_test.go`
- `backend/internal/adapter/external/consultahistorica/consultador.go`
- `backend/internal/adapter/external/consultahistorica/consultador_test.go`
- `backend/internal/adapter/rest/dto/consulta_facturas.go`
- `backend/internal/adapter/rest/dto/importacion_facturas.go`
- `backend/internal/adapter/rest/dto/respuestas_facturacion_web.go`
- `backend/internal/adapter/rest/errors/mapeador.go`
- `backend/internal/adapter/rest/handler/consulta_facturas_handler.go`
- `backend/internal/adapter/rest/handler/facturacion_web_handler.go`
- `backend/internal/adapter/rest/handler/handlers.go`
- `backend/internal/adapter/rest/handler/importacion_facturas_handler.go`
- `backend/internal/adapter/rest/router/arca.go`
- `backend/internal/core/domain/facturacion.go`
- `backend/internal/core/domain/importacion_factura.go`
- `backend/internal/core/domain/importacion_validacion.go`
- `backend/internal/core/domain/importacion_validacion_test.go`
- `backend/internal/core/domain/integridad_fiscal.go`
- `backend/internal/core/port/in/casos_uso.go`
- `backend/internal/core/port/in/importacion_facturas.go`
- `backend/internal/core/port/out/importacion_facturas.go`
- `backend/internal/core/port/out/repositorios.go`
- `backend/internal/usecase/comprobante/generar_comprobante.go`
- `backend/internal/usecase/importarfacturas/caso_uso.go`
- `backend/internal/usecase/importarfacturas/caso_uso_test.go`

## Inventario SDK

- `README.md`
- `contrato_importacion_test.go`
- `docs/CONTRATO.md`
- `docs/arquitectura.md`
- `docs/entrega-importacion-fiscal.md`
- `docs/folder-structure.md`
- `docs/importacion-facturas.md`
- `docs/openapi.yaml`
- `docs/validacion.md`
- `examples/importar-factura/main.go`
- `importacion_facturas.go`
- `importacion_facturas_test.go`
- `integration/consumer/importacion_test.go`
- `modelos/facturacion/importacion.go`
- `modelos/facturacion/persistidas.go`
