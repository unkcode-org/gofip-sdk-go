// Package gofip implementa un cliente idiomático para la API v2 de GOFIP.
//
// NuevoCliente usa Bearer para emisión, consultas ARCA, documentos y consultas
// persistidas por CAE/documento ID. Estas últimas leen MySQL a nivel de emisor.
// NuevoClienteSEP usa IntegrationKey exclusivamente para las operaciones M2M
// de creación y recuperación de acceso Embedded; no implementa el navegador.
package gofip
