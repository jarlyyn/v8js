package httpv8

import (
	"net/url"

	"github.com/herb-go/herbplugin"
	"github.com/herb-go/plugins/addons/httpaddon"
	v8js "github.com/jarlyyn/v8js"
)

type Builder func(r *v8js.Context, req *Request) *v8js.JsValue

var DefaultBuilder = func(r *v8js.Context, req *Request) *v8js.JsValue {
	obj := r.NewObject()
	obj.SetObjectMethod(r, "GetID", req.GetID)
	obj.SetObjectMethod(r, "GetURL", req.GetURL)
	obj.SetObjectMethod(r, "SetURL", req.SetURL)
	obj.SetObjectMethod(r, "GetProxy", req.GetProxy)
	obj.SetObjectMethod(r, "SetProxy", req.SetProxy)
	obj.SetObjectMethod(r, "GetMethod", req.GetMethod)
	obj.SetObjectMethod(r, "SetMethod", req.SetMethod)
	obj.SetObjectMethod(r, "GetBody", req.GetBody)
	obj.SetObjectMethod(r, "SetBody", req.SetBody)
	obj.SetObjectMethod(r, "FinishedAt", req.FinishedAt)
	obj.SetObjectMethod(r, "ExecuteStatus", req.ExecuteStatus)
	obj.SetObjectMethod(r, "ResetHeader", req.ResetHeader)
	obj.SetObjectMethod(r, "SetHeader", req.SetHeader)
	obj.SetObjectMethod(r, "AddHeader", req.AddHeader)
	obj.SetObjectMethod(r, "DelHeader", req.DelHeader)
	obj.SetObjectMethod(r, "GetHeader", req.GetHeader)
	obj.SetObjectMethod(r, "HeaderValues", req.HeaderValues)
	obj.SetObjectMethod(r, "HeaderFields", req.HeaderFields)
	obj.SetObjectMethod(r, "ResponseStatusCode", req.ResponseStatusCode)
	obj.SetObjectMethod(r, "ResponseBody", req.ResponseBody)
	obj.SetObjectMethod(r, "ResponseHeader", req.ResponseHeader)
	obj.SetObjectMethod(r, "ResponseHeaderValues", req.ResponseHeaderValues)
	obj.SetObjectMethod(r, "ResponseHeaderFields", req.ResponseHeaderFields)
	obj.SetObjectMethod(r, "Execute", req.Execute)
	return obj
}

type Request struct {
	Request *httpaddon.Request
}

func (req *Request) GetID(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	return call.Context().NewString(req.Request.GetID())
}

func (req *Request) GetURL(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	return call.Context().NewString(req.Request.GetURL())
}
func (req *Request) SetURL(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	req.Request.SetURL(call.GetArg(0).String())
	return nil
}
func (req *Request) GetProxy(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	return call.Context().NewString(req.Request.GetProxy())
}
func (req *Request) SetProxy(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	req.Request.SetProxy(call.GetArg(0).String())
	return nil
}

func (req *Request) GetMethod(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	return call.Context().NewString(req.Request.GetMethod())
}
func (req *Request) SetMethod(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	req.Request.SetMethod(call.GetArg(0).String())
	return nil
}
func (req *Request) GetBody(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	return call.Context().NewString(string(req.Request.GetBody()))
}
func (req *Request) SetBody(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	req.Request.SetBody([]byte(call.GetArg(0).String()))
	return nil
}

func (req *Request) FinishedAt(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	return call.Context().NewInt64(req.Request.FinishedAt())

}
func (req *Request) ExecuteStatus(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	return call.Context().NewInt32(int32(req.Request.ExecuteStauts()))
}
func (req *Request) ResetHeader(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	req.Request.ResetHeader()
	return nil
}
func (req *Request) SetHeader(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	req.Request.SetHeader(call.GetArg(0).String(), call.GetArg(1).String())
	return nil
}
func (req *Request) AddHeader(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	req.Request.AddHeader(call.GetArg(0).String(), call.GetArg(1).String())
	return nil
}
func (req *Request) DelHeader(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	req.Request.DelHeader(call.GetArg(0).String())
	return nil

}
func (req *Request) GetHeader(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	return call.Context().NewString(req.Request.GetHeader(call.GetArg(0).String()))

}
func (req *Request) HeaderValues(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	result := req.Request.HeaderValues(call.GetArg(0).String())
	var output = make([]*v8js.JsValue, len(result))
	for i, v := range result {
		output[i] = call.Context().NewString(v)
	}
	return call.Context().NewArray(output...)
}
func (req *Request) HeaderFields(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	result := req.Request.HeaderFields()
	var output = make([]*v8js.JsValue, len(result))
	for i, v := range result {
		output[i] = call.Context().NewString(v)
	}
	return call.Context().NewArray(output...)

}

func (req *Request) ResponseStatusCode(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	return call.Context().NewInt32(int32(req.Request.ResponseStatusCode()))
}
func (req *Request) ResponseBody(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	return call.Context().NewString(string(req.Request.ResponseBody()))
}
func (req *Request) ResponseHeader(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	return call.Context().NewString(req.Request.ResponseHeader(call.GetArg(0).String()))

}
func (req *Request) ResponseHeaderValues(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	result := req.Request.ResponseHeaderValues(call.GetArg(0).String())

	var output = make([]*v8js.JsValue, len(result))
	for i, v := range result {
		output[i] = call.Context().NewString(v)
	}
	return call.Context().NewArray(output...)
}
func (req *Request) ResponseHeaderFields(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	result := req.Request.ResponseHeaderFields()

	var output = make([]*v8js.JsValue, len(result))
	for i, v := range result {
		output[i] = call.Context().NewString(v)
	}
	return call.Context().NewArray(output...)

}
func (req *Request) Execute(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	req.Request.MustExecute()
	return nil
}

type Addon struct {
	Addon   *httpaddon.Addon
	Builder Builder
}

func (a *Addon) ParseURL(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	rawurl := call.Args()[0].String()
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil
	}
	result := call.Context().NewObject()
	result.Set("Host", call.Context().NewString(u.Host))
	result.Set("Hostname", call.Context().NewString(u.Host))
	result.Set("Scheme", call.Context().NewString(u.Scheme))
	result.Set("Path", call.Context().NewString(u.Path))
	result.Set("Query", call.Context().NewString(u.RawQuery))
	result.Set("User", call.Context().NewString(u.User.Username()))
	p, _ := u.User.Password()
	result.Set("Password", call.Context().NewString(p))
	result.Set("Port", call.Context().NewString(u.Port()))
	result.Set("Fragment", call.Context().NewString(u.Fragment))
	return result
}
func (a *Addon) NewRequest(call *v8js.FunctionCallbackInfo) *v8js.JsValue {

	method := call.GetArg(0).String()
	url := call.GetArg(1).String()
	req := a.Addon.Create(method, url)
	return a.Builder(call.Context(), &Request{req})
}

func (a *Addon) Convert(r *v8js.Context) *v8js.JsValue {
	obj := r.NewObject()
	obj.SetObjectMethod(r, "New", a.NewRequest)
	obj.SetObjectMethod(r, "ParseURL", a.ParseURL)

	return obj
}
func Create(p herbplugin.Plugin) *Addon {
	return &Addon{
		Addon:   httpaddon.Create(p),
		Builder: DefaultBuilder,
	}
}
