# Arquitectura de integración GOFIP V2

El integrador mantiene dos clientes independientes. `NuevoCliente` conserva el
transporte Bearer y todos los métodos públicos anteriores. `NuevoClienteSEP`
construye IntegrationKey y políticas de seguridad propias; no reutiliza la
construcción de autenticación Bearer. Se comparten el envoltorio JSON `data` y el
parser de errores; SEP sanea el resultado antes de devolverlo. No se agregan
dependencias ni se modifica la semántica del transporte anterior.

| Operación | Fuente | Cliente |
| --- | --- | --- |
| EmitirFactura | GOFIP → ARCA, persistencia | Bearer |
| ObtenerComprobante | Consulta directa a ARCA | Bearer |
| ObtenerFacturaPorCAE / ObtenerFacturaPorId | MySQL exclusivamente | Bearer |
| Documentos | GOFIP y almacenamiento de documentos | Bearer |
| CrearSesion / CrearBootstrap | GOFIP Embedded M2M | SEP |

La ampliación de lectura se aplica sólo a los GET nuevos: una API key permite
consultar todos los puntos de venta del emisor asociado a su punto de venta.
La emisión continúa acotada al punto de venta de esa credencial. El consumidor
no puede proporcionar un emisor para ampliar el acceso.

Los modelos fiscales persistidos tienen identidad distinta del comprobante ARCA:
incluyen documento/factura ID, estado GOFIP, ítems guardados y metadata PDF. Las
consultas por CAE e ID comparten `RespuestaFacturaPersistida`. No se recuperan datos
desde Legacy, ARCA o entidades mutables como fallback. No existe filtro de ambiente
porque el modelo actual no persiste ese dato. No se agregan migraciones ni índices
únicos de CAE. Un CAE ambiguo produce 409 y admite filtros explícitos.

SEP nunca guarda cookies de navegador ni hace intercambio de bootstrap. Sus
respuestas contienen capacidades de corta duración con acceso explícito mediante
`Secreto.Valor()`. JSON de diagnóstico, fmt y slog redactan tanto tokens como URL de
documento. El integrador controla la entrega deliberada al frontend y debe evitar
que su RoundTripper registre las credenciales o cuerpos observados.

El contrato canónico se mantiene en `backend/docs/openapi.yaml` de gofip-backend;
`docs/openapi.yaml` es una copia exacta. Cada cambio debe comenzar en el backend,
incrementar su versión y changelog y luego sincronizar esta copia. El SDK conserva
el módulo `/v2` y Go mínimo 1.22.2.
