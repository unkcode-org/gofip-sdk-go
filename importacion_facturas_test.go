package gofip

import (
 "context"
 "errors"
 "fmt"
 "io"
 "net/http"
 "net/http/httptest"
 "strings"
 "testing"

 "github.com/unkcode-org/gofip-sdk-go/v2/modelos/facturacion"
)
func solicitudHistoricaSDK()facturacion.SolicitudImportarFactura{return facturacion.SolicitudImportarFactura{Factura:facturacion.FacturaHistorica{PuntoVenta:3,TipoComprobante:6,NumeroComprobante:12345,CAE:"00123456789012"},Procedencia:facturacion.ProcedenciaFactura{Sistema:"gofip_legacy"}}}
const respuestaHistoricaSDK=`{"documento_id":"01900000-0000-7000-8000-000000000001","importado":true,"ya_existia":false,"verificacion":{"metodo":"wsfe.FECompConsultar","verificado_at":"2026-10-16T20:00:00Z"},"documento":{"documento_id":"01900000-0000-7000-8000-000000000001","emisor_id":"emisor-a","punto_venta_id":"pv-3","estado":"autorizado","origen":"legacy_import","factura":{"cae":"00123456789012","historico":{"cotizacion":"1234.123456","importe_total":"121.00"}},"pdf":{"disponible":false}}}`
func TestImportacionSDKContratoTipadoSinSeleccionAmbiente(t *testing.T){
 calls:=0
 s:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  calls++;if r.Method!="POST"||r.URL.Path!="/api/v2/facturas/importaciones"||r.Header.Get("Authorization")!="Bearer clave-prueba"||r.Header.Get("Idempotency-Key")!="historico:123"||r.Header.Get("X-External-Reference")!="cliente:123"||r.Header.Get("X-ARCA-Ambiente")!="" {t.Error("contrato/autenticación/ambiente incorrecto")}
  b,_:=io.ReadAll(r.Body);if !strings.Contains(string(b),`"data":{"factura":` )||strings.Contains(string(b),"ambiente")||strings.Contains(string(b),"emisor_id"){t.Error("body no tipado o amplía autoridad")}
  w.Header().Set("Content-Type","application/json");w.WriteHeader(201);_,_=io.WriteString(w,respuestaHistoricaSDK)
 }));defer s.Close()
 c,err:=NuevoCliente(OpcionesCliente{URLBase:s.URL,ClaveAPI:"clave-prueba",AmbienteARCA:AmbienteHomologacion,ClienteHTTP:s.Client()});if err!=nil{t.Fatal(err)}
 r,err:=c.Facturacion.ImportarFactura(context.Background(),solicitudHistoricaSDK(),OpcionesImportacionFactura{ClaveIdempotencia:"historico:123",ReferenciaExterna:"cliente:123"});if err!=nil{t.Fatal(err)}
 if calls!=1||r.DocumentoID!=documentoFiscalPrueba||r.Documento.Factura.Historico.Cotizacion==nil||*r.Documento.Factura.Historico.Cotizacion!="1234.123456"{t.Fatal("datos históricos alterados")}
}
func TestImportacionSDKErroresSegurosLimitesYNoRetries(t *testing.T){
 secret:="clave-api-no-exponer"
 for _,tc:=range []struct{name string;status int;body string}{{"conflicto",409,`{"code":"13101","detail":"`+secret+`"}`},{"no-json",502,secret},{"body-grande",503,strings.Repeat("x",(2<<20)+1)},{"respuesta-invalida",201,`{"documento_id":"`+secret+`"}`},{"json-adicional",201,respuestaHistoricaSDK+` {}`}}{t.Run(tc.name,func(t *testing.T){
  calls:=0;transport:=roundTripperPrueba(func(r *http.Request)(*http.Response,error){calls++;return &http.Response{StatusCode:tc.status,Header:make(http.Header),Body:io.NopCloser(strings.NewReader(tc.body)),Request:r},nil})
  c,_:=NuevoCliente(OpcionesCliente{URLBase:"https://example.test",ClaveAPI:secret,ClienteHTTP:&http.Client{Transport:transport}})
  _,err:=c.Facturacion.ImportarFactura(context.Background(),solicitudHistoricaSDK(),OpcionesImportacionFactura{ClaveIdempotencia:"uno"})
  if err==nil||strings.Contains(fmt.Sprintf("%+v",err),secret)||calls!=1{t.Fatalf("error/retry inseguro: %v calls=%d",err,calls)}
  if tc.status==409{var api *ErrorAPI;if !errors.As(err,&api)||api.Estado!=409||api.Codigo!="13101"{t.Fatal("perdió error tipado",err)}}
 })}
 transport:=roundTripperPrueba(func(*http.Request)(*http.Response,error){return nil,errors.New(secret)})
 c,_:=NuevoCliente(OpcionesCliente{URLBase:"https://example.test",ClaveAPI:secret,ClienteHTTP:&http.Client{Transport:transport}})
 _,err:=c.Facturacion.ImportarFactura(context.Background(),solicitudHistoricaSDK(),OpcionesImportacionFactura{ClaveIdempotencia:"uno"});if err==nil||strings.Contains(fmt.Sprintf("%+v",err),secret){t.Fatal("filtró error de transporte")}
}
func TestImportacionSDKCancelacionYRedireccion(t *testing.T){
 calls:=0
 s:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){calls++;http.Redirect(w,r,"/otra-ruta",http.StatusTemporaryRedirect)}));defer s.Close()
 c:=nuevoClientePrueba(t,s.URL,s.Client());ctx,cancel:=context.WithCancel(context.Background());cancel()
 _,err:=c.Facturacion.ImportarFactura(ctx,solicitudHistoricaSDK(),OpcionesImportacionFactura{ClaveIdempotencia:"uno"});if !errors.Is(err,context.Canceled)||calls!=0{t.Fatal("cancelación perdida",err)}
 _,err=c.Facturacion.ImportarFactura(context.Background(),solicitudHistoricaSDK(),OpcionesImportacionFactura{ClaveIdempotencia:"uno"});if err==nil||calls!=1{t.Fatal("siguió redirect de importación",err)}
}
