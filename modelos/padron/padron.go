package padron

import "time"

type TipoDocumentoP13 string

const (
	TipoDocumentoDNI  TipoDocumentoP13 = "dni"
	TipoDocumentoCUIT TipoDocumentoP13 = "cuit"
	TipoDocumentoCUIL TipoDocumentoP13 = "cuil"
)

type RespuestaPersonaPadron struct {
	PersonaP4  *PersonaP4  `json:"persona_p4,omitempty"`
	PersonaP5  *PersonaP5  `json:"persona_p5,omitempty"`
	PersonaP10 *PersonaP10 `json:"persona_p10,omitempty"`
	PersonaP13 *PersonaP13 `json:"persona_p13,omitempty"`
}

type PersonaP4 struct {
	Actividad                 []ActividadPadron   `json:"actividad,omitempty"`
	Apellido                  *string             `json:"apellido,omitempty"`
	CantidadSociosEmpresaMono *int                `json:"cantidad_socios_empresa_mono,omitempty"`
	Categoria                 []Categoria         `json:"categoria,omitempty"`
	ClaveInactivaAsociada     []int64             `json:"clave_inactiva_asociada,omitempty"`
	Dependencia               *Dependencia        `json:"dependencia,omitempty"`
	Domicilio                 []Domicilio         `json:"domicilio,omitempty"`
	Email                     []EmailPadron       `json:"email,omitempty"`
	EstadoClave               *string             `json:"estado_clave,omitempty"`
	FechaContratoSocial       *time.Time          `json:"fecha_contrato_social,omitempty"`
	FechaFallecimiento        *time.Time          `json:"fecha_fallecimiento,omitempty"`
	FechaInscripcion          *time.Time          `json:"fecha_inscripcion,omitempty"`
	FechaJubilado             *time.Time          `json:"fecha_jubilado,omitempty"`
	FechaNacimiento           *time.Time          `json:"fecha_nacimiento,omitempty"`
	FechaVencimientoMigracion *time.Time          `json:"fecha_vencimiento_migracion,omitempty"`
	FormaJuridica             *string             `json:"forma_juridica,omitempty"`
	IDPersona                 *int64              `json:"id_persona,omitempty"`
	Impuesto                  []ImpuestoDetallado `json:"impuesto,omitempty"`
	LeyJubilacion             *int                `json:"ley_jubilacion,omitempty"`
	LocalidadInscripcion      *string             `json:"localidad_inscripcion,omitempty"`
	MesCierre                 *int                `json:"mes_cierre,omitempty"`
	Nombre                    *string             `json:"nombre,omitempty"`
	NumeroDocumento           *string             `json:"numero_documento,omitempty"`
	NumeroInscripcion         *int64              `json:"numero_inscripcion,omitempty"`
	OrganismoInscripcion      *string             `json:"organismo_inscripcion,omitempty"`
	OrganismoOriginante       *string             `json:"organismo_originante,omitempty"`
	PorcentajeCapitalNacional *float64            `json:"porcentaje_capital_nacional,omitempty"`
	ProvinciaInscripcion      *string             `json:"provincia_inscripcion,omitempty"`
	RazonSocial               *string             `json:"razon_social,omitempty"`
	Regimen                   []RegimenExtendido  `json:"regimen,omitempty"`
	Relacion                  []Relacion          `json:"relacion,omitempty"`
	Sexo                      *string             `json:"sexo,omitempty"`
	Telefono                  []Telefono          `json:"telefono,omitempty"`
	TipoClave                 *string             `json:"tipo_clave,omitempty"`
	TipoDocumento             *string             `json:"tipo_documento,omitempty"`
	TipoOrganismoOriginante   *string             `json:"tipo_organismo_originante,omitempty"`
	TipoPersona               *string             `json:"tipo_persona,omitempty"`
	TipoResidencia            *string             `json:"tipo_residencia,omitempty"`
}

type PersonaP5 struct {
	DatosGenerales      *DatosGenerales      `json:"datos_generales,omitempty"`
	DatosMonotributo    *DatosMonotributo    `json:"datos_monotributo,omitempty"`
	DatosRegimenGeneral *DatosRegimenGeneral `json:"datos_regimen_general,omitempty"`
	ErrorConstancia     *ErrorConstancia     `json:"error_constancia,omitempty"`
	ErrorMonotributo    *ErrorMonotributo    `json:"error_monotributo,omitempty"`
	ErrorRegimenGeneral *ErrorRegimenGeneral `json:"error_regimen_general,omitempty"`
}

type PersonaP10 struct {
	Apellido                      *string      `json:"apellido,omitempty"`
	ClaveInactivaAsociada         []int64      `json:"clave_inactiva_asociada,omitempty"`
	Dependencia                   *Dependencia `json:"dependencia,omitempty"`
	DescripcionActividadPrincipal *string      `json:"descripcion_actividad_principal,omitempty"`
	Domicilio                     []Domicilio  `json:"domicilio,omitempty"`
	EstadoClave                   *string      `json:"estado_clave,omitempty"`
	IDActividadPrincipal          *int64       `json:"id_actividad_principal,omitempty"`
	IDPersona                     *int64       `json:"id_persona,omitempty"`
	Nombre                        *string      `json:"nombre,omitempty"`
	NumeroDocumento               *string      `json:"numero_documento,omitempty"`
	RazonSocial                   *string      `json:"razon_social,omitempty"`
	TipoClave                     *string      `json:"tipo_clave,omitempty"`
	TipoDocumento                 *string      `json:"tipo_documento,omitempty"`
	TipoPersona                   *string      `json:"tipo_persona,omitempty"`
}

type PersonaP13 struct {
	Apellido                      *string              `json:"apellido,omitempty"`
	ClaveInactivaAsociada         []int64              `json:"clave_inactiva_asociada,omitempty"`
	DescripcionActividadPrincipal *string              `json:"descripcion_actividad_principal,omitempty"`
	Domicilio                     []DomicilioExtendido `json:"domicilio,omitempty"`
	EstadoClave                   *string              `json:"estado_clave,omitempty"`
	FechaContratoSocial           *time.Time           `json:"fecha_contrato_social,omitempty"`
	FechaFallecimiento            *time.Time           `json:"fecha_fallecimiento,omitempty"`
	FechaNacimiento               *time.Time           `json:"fecha_nacimiento,omitempty"`
	FormaJuridica                 *string              `json:"forma_juridica,omitempty"`
	IDActividadPrincipal          *int64               `json:"id_actividad_principal,omitempty"`
	IDPersona                     *int64               `json:"id_persona,omitempty"`
	MesCierre                     *int                 `json:"mes_cierre,omitempty"`
	Nombre                        *string              `json:"nombre,omitempty"`
	NumeroDocumento               *string              `json:"numero_documento,omitempty"`
	PeriodoActividadPrincipal     *int                 `json:"periodo_actividad_principal,omitempty"`
	RazonSocial                   *string              `json:"razon_social,omitempty"`
	TipoClave                     *string              `json:"tipo_clave,omitempty"`
	TipoDocumento                 *string              `json:"tipo_documento,omitempty"`
	TipoPersona                   *string              `json:"tipo_persona,omitempty"`
}

type DatosGenerales struct {
	Apellido            *string               `json:"apellido,omitempty"`
	Caracterizacion     []Caracterizacion     `json:"caracterizacion,omitempty"`
	Dependencia         *DependenciaDetallada `json:"dependencia,omitempty"`
	DomicilioFiscal     *Domicilio            `json:"domicilio_fiscal,omitempty"`
	EsSucesion          *string               `json:"es_sucesion,omitempty"`
	EstadoClave         *string               `json:"estado_clave,omitempty"`
	FechaContratoSocial *time.Time            `json:"fecha_contrato_social,omitempty"`
	FechaFallecimiento  *time.Time            `json:"fecha_fallecimiento,omitempty"`
	IDPersona           *int64                `json:"id_persona,omitempty"`
	MesCierre           *int                  `json:"mes_cierre,omitempty"`
	Nombre              *string               `json:"nombre,omitempty"`
	RazonSocial         *string               `json:"razon_social,omitempty"`
	TipoClave           *string               `json:"tipo_clave,omitempty"`
	TipoPersona         *string               `json:"tipo_persona,omitempty"`
}

type DatosMonotributo struct {
	Actividad               []ActividadPadron `json:"actividad,omitempty"`
	ActividadMonotributista *ActividadPadron  `json:"actividad_monotributista,omitempty"`
	CategoriaMonotributo    *CategoriaFiscal  `json:"categoria_monotributo,omitempty"`
	ComponenteDeSociedad    []RelacionFiscal  `json:"componente_de_sociedad,omitempty"`
	Impuesto                []Impuesto        `json:"impuesto,omitempty"`
}

type DatosRegimenGeneral struct {
	Actividad         []ActividadPadron `json:"actividad,omitempty"`
	CategoriaAutonomo *CategoriaFiscal  `json:"categoria_autonomo,omitempty"`
	Impuesto          []Impuesto        `json:"impuesto,omitempty"`
	Regimen           []Regimen         `json:"regimen,omitempty"`
}

type ErrorConstancia struct {
	Apellido  *string  `json:"apellido,omitempty"`
	Error     []string `json:"error,omitempty"`
	IDPersona *int64   `json:"id_persona,omitempty"`
	Nombre    *string  `json:"nombre,omitempty"`
}

type ErrorMonotributo struct {
	Error   []string `json:"error,omitempty"`
	Mensaje *string  `json:"mensaje,omitempty"`
}

type ErrorRegimenGeneral struct {
	Error   []string `json:"error,omitempty"`
	Mensaje *string  `json:"mensaje,omitempty"`
}

type Caracterizacion struct {
	DescripcionCaracterizacion *string `json:"descripcion_caracterizacion,omitempty"`
	IDCaracterizacion          *int    `json:"id_caracterizacion,omitempty"`
	Periodo                    *int    `json:"periodo,omitempty"`
}

type Telefono struct {
	Numero       *int64  `json:"numero,omitempty"`
	TipoLinea    *string `json:"tipo_linea,omitempty"`
	TipoTelefono *string `json:"tipo_telefono,omitempty"`
}

type Impuesto struct {
	DescripcionImpuesto *string `json:"descripcion_impuesto,omitempty"`
	EstadoImpuesto      *string `json:"estado_impuesto,omitempty"`
	IDImpuesto          *int    `json:"id_impuesto,omitempty"`
	Motivo              *string `json:"motivo,omitempty"`
	Periodo             *int    `json:"periodo,omitempty"`
}

type ImpuestoDetallado struct {
	DescripcionImpuesto *string    `json:"descripcion_impuesto,omitempty"`
	DiaPeriodo          *int       `json:"dia_periodo,omitempty"`
	Estado              *string    `json:"estado,omitempty"`
	FechaInscripcion    *time.Time `json:"ff_inscripcion,omitempty"`
	IDImpuesto          *int       `json:"id_impuesto,omitempty"`
	Periodo             *int       `json:"periodo,omitempty"`
}

type Regimen struct {
	DescripcionRegimen *string `json:"descripcion_regimen,omitempty"`
	IDImpuesto         *int    `json:"id_impuesto,omitempty"`
	IDRegimen          *int    `json:"id_regimen,omitempty"`
	Periodo            *int    `json:"periodo,omitempty"`
	TipoRegimen        *string `json:"tipo_regimen,omitempty"`
}

type RegimenExtendido struct {
	DescripcionRegimen *string `json:"descripcion_regimen,omitempty"`
	DiaPeriodo         *int    `json:"dia_periodo,omitempty"`
	Estado             *string `json:"estado,omitempty"`
	IDImpuesto         *int    `json:"id_impuesto,omitempty"`
	IDRegimen          *int    `json:"id_regimen,omitempty"`
	Periodo            *int    `json:"periodo,omitempty"`
	TipoRegimen        *string `json:"tipo_regimen,omitempty"`
}

type RelacionFiscal struct {
	ApellidoPersonaAsociada    *string    `json:"apellido_persona_asociada,omitempty"`
	FechaRelacion              *time.Time `json:"ff_relacion,omitempty"`
	FechaVencimiento           *time.Time `json:"ff_vencimiento,omitempty"`
	IDPersonaAsociada          *int64     `json:"id_persona_asociada,omitempty"`
	NombrePersonaAsociada      *string    `json:"nombre_persona_asociada,omitempty"`
	RazonSocialPersonaAsociada *string    `json:"razon_social_persona_asociada,omitempty"`
	TipoComponente             *string    `json:"tipo_componente,omitempty"`
}

type Relacion struct {
	FechaRelacion     *time.Time `json:"ff_relacion,omitempty"`
	FechaVencimiento  *time.Time `json:"ff_vencimiento,omitempty"`
	IDPersona         *int64     `json:"id_persona,omitempty"`
	IDPersonaAsociada *int64     `json:"id_persona_asociada,omitempty"`
	SubtipoRelacion   *string    `json:"subtipo_relacion,omitempty"`
	TipoRelacion      *string    `json:"tipo_relacion,omitempty"`
}

type CategoriaFiscal struct {
	DescripcionCategoria *string `json:"descripcion_categoria,omitempty"`
	IDCategoria          *int    `json:"id_categoria,omitempty"`
	IDImpuesto           *int    `json:"id_impuesto,omitempty"`
	Periodo              *int    `json:"periodo,omitempty"`
}

type Categoria struct {
	DescripcionCategoria *string `json:"descripcion_categoria,omitempty"`
	Estado               *string `json:"estado,omitempty"`
	IDCategoria          *int    `json:"id_categoria,omitempty"`
	IDImpuesto           *int    `json:"id_impuesto,omitempty"`
	Periodo              *int    `json:"periodo,omitempty"`
}

type EmailPadron struct {
	Direccion *string `json:"direccion,omitempty"`
	Estado    *string `json:"estado,omitempty"`
	TipoEmail *string `json:"tipo_email,omitempty"`
}

type ActividadPadron struct {
	DescripcionActividad *string `json:"descripcion_actividad,omitempty"`
	IDActividad          *int64  `json:"id_actividad,omitempty"`
	Nomenclador          *int    `json:"nomenclador,omitempty"`
	Orden                *int    `json:"orden,omitempty"`
	Periodo              *int    `json:"periodo,omitempty"`
}

type Dependencia struct {
	DescripcionDependencia *string `json:"descripcion_dependencia,omitempty"`
	IDDependencia          *int    `json:"id_dependencia,omitempty"`
}

type DependenciaDetallada struct {
	CodigoPostal           *string `json:"cod_postal,omitempty"`
	DescripcionDependencia *string `json:"descripcion_dependencia,omitempty"`
	DescripcionProvincia   *string `json:"descripcion_provincia,omitempty"`
	Direccion              *string `json:"direccion,omitempty"`
	IDDependencia          *int    `json:"id_dependencia,omitempty"`
	IDProvincia            *int    `json:"id_provincia,omitempty"`
	Localidad              *string `json:"localidad,omitempty"`
}

type Domicilio struct {
	CodigoPostal         *string `json:"cod_postal,omitempty"`
	DatoAdicional        *string `json:"dato_adicional,omitempty"`
	DescripcionProvincia *string `json:"descripcion_provincia,omitempty"`
	Direccion            *string `json:"direccion,omitempty"`
	IDProvincia          *int    `json:"id_provincia,omitempty"`
	Localidad            *string `json:"localidad,omitempty"`
	TipoDatoAdicional    *string `json:"tipo_dato_adicional,omitempty"`
	TipoDomicilio        *string `json:"tipo_domicilio,omitempty"`
}

type DomicilioExtendido struct {
	Calle                *string `json:"calle,omitempty"`
	CodigoPostal         *string `json:"codigo_postal,omitempty"`
	DatoAdicional        *string `json:"dato_adicional,omitempty"`
	DescripcionProvincia *string `json:"descripcion_provincia,omitempty"`
	Direccion            *string `json:"direccion,omitempty"`
	EstadoDomicilio      *string `json:"estado_domicilio,omitempty"`
	IDProvincia          *int    `json:"id_provincia,omitempty"`
	Localidad            *string `json:"localidad,omitempty"`
	Manzana              *string `json:"manzana,omitempty"`
	Numero               *int    `json:"numero,omitempty"`
	OficinaDptoLocal     *string `json:"oficina_dpto_local,omitempty"`
	Piso                 *string `json:"piso,omitempty"`
	Sector               *string `json:"sector,omitempty"`
	TipoDatoAdicional    *string `json:"tipo_dato_adicional,omitempty"`
	TipoDomicilio        *string `json:"tipo_domicilio,omitempty"`
	Torre                *string `json:"torre,omitempty"`
}
