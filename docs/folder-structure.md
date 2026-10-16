# Estructura del SDK

```text
gofip-sdk-go/
├── cliente.go
├── transporte.go
├── utilidades.go
├── catalogos.go
├── padron.go
├── facturacion.go
├── documentos.go
├── modelos/
│   ├── catalogos/
│   ├── documentos/
│   ├── facturacion/
│   └── padron/
├── examples/
├── docs/
│   └── openapi.yaml
├── go.mod
└── README.md
```

El paquete raíz `gofip` contiene el cliente, los servicios, el transporte y los
errores compartidos. Los modelos públicos grandes se separan por dominio para
evitar un paquete raíz difícil de navegar.

`docs/openapi.yaml` es la fuente de verdad del contrato HTTP.

El cliente SEP está en `cliente_sep.go`, `sep.go` y `transporte_sep.go`; sus modelos están en `modelos/sep`. Las consultas persistidas están en `facturas_persistidas.go` y `modelos/facturacion/persistidas.go`. Ver [arquitectura](arquitectura.md).
