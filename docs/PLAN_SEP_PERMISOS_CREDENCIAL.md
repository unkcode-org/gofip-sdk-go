# Cambio SEP propuesto: permisos determinados por Integration Key

Fecha: 2026-10-18. Estado: auditoría completa; implementación pendiente de aprobación TDD coordinada con frontend.

Rama: `codex/sep-permisos-por-credencial`, creada desde `origin/main` actualizado (`bc407552513560fa921b4bc637bac416239f2e33`); no existe develop en este repositorio.

## Hallazgo

El SDK sólo expone SEP M2M mediante CrearSesion/CrearBootstrap. Ambas funciones reciben SolicitudCrearSesionEmbedded, validan Permisos y serializan `{"data":{"permissions":[...]}}`. No implementa administración SEP ni canje de navegador. El catálogo tiene doce permisos, incluido invoices:create; no existe invoices:write.

## API propuesta

```go
sesion, err := clienteSEP.CrearSesion(ctx)
bootstrap, err := clienteSEP.CrearBootstrap(ctx)
```

Eliminar DTO público vacío y validación de permisos del consumidor. El transporte enviará POST sin cuerpo; la clave configurada determina integración, emisor, credencial y todos sus permisos. Mantener tipos legítimos del catálogo y toda la seguridad de capacidades, redacción, no redirects, no cookies y no retry de POST. No agregar administración al SDK.

Actualizar modelos/cliente/transporte/tests/ejemplos/README, consumidor externo anidado y procedencia del OpenAPI. Contrato propuesto 4.0.0, mismo /api/v2, YAML idéntico byte a byte al backend/frontend. El YAML actual sigue en 3.0.0, SHA-256 b9580c0245564ce8cf3d11930b36650b8b6e8934b3910c492b4c9a5ae0557238.

La auditoría completa y el plan TDD están en backend/docs/AUDITORIA_SEP_PERMISOS_CREDENCIAL.md y backend/docs/PLAN_TDD_SEP_PERMISOS_CREDENCIAL.md del repositorio canónico gofip-backend. El frontend guarda también el plan en .harness/tasks/sep-permisos-por-credencial/tdd-test-plan.md. Se requiere aprobación del plan antes de implementación.

## Versionado y Elixir

Propuesta de siguiente prerelease: v2.0.0-rc.4. Es incompatible a nivel de API Go con v2.0.0-rc.3, dentro de una serie prerelease; no crear tag ni release. Un SDK major estable requeriría respetar SemVer y evaluar v3 si correspondiera.

Elixir deberá actualizar posteriormente el SDK y quitar SolicitudCrearSesionEmbedded{Permisos: ...} en ambas llamadas. GOFIP_SEP_INTEGRATION_KEY basta; no agregar env vars para permisos. No se modifica Elixir en esta tarea.

## Verificación previa

`go test ./...` pasó sobre la base oficial. No incluye el módulo anidado integration/consumer. Race/vet, módulo externo, pruebas del contrato nuevo y validación completa se ejecutarán después de implementar. No se publicó ningún cambio.
