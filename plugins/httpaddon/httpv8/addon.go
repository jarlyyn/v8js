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

func (req *Request) GetID(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	return call.Context().NewString(req.Request.GetID()).Consume()
}

func (req *Request) GetURL(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	return call.Context().NewString(req.Request.GetURL()).Consume()
}
func (req *Request) SetURL(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	req.Request.SetURL(call.GetArg(0).String())
	return nil
}
func (req *Request) GetProxy(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	return call.Context().NewString(req.Request.GetProxy()).Consume()
}
func (req *Request) SetProxy(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	req.Request.SetProxy(call.GetArg(0).String())
	return nil
}

func (req *Request) GetMethod(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	return call.Context().NewString(req.Request.GetMethod()).Consume()
}
func (req *Request) SetMethod(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	req.Request.SetMethod(call.GetArg(0).String())
	return nil
}
func (req *Request) GetBody(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	return call.Context().NewString(string(req.Request.GetBody())).Consume()
}
func (req *Request) SetBody(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	req.Request.SetBody([]byte(call.GetArg(0).String()))
	return nil
}

func (req *Request) FinishedAt(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	return call.Context().NewInt64(req.Request.FinishedAt()).Consume()

}
func (req *Request) ExecuteStatus(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	return call.Context().NewInt32(int32(req.Request.ExecuteStauts())).Consume()
}
func (req *Request) ResetHeader(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	req.Request.ResetHeader()
	return nil
}
func (req *Request) SetHeader(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	req.Request.SetHeader(call.GetArg(0).String(), call.GetArg(1).String())
	return nil
}
func (req *Request) AddHeader(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	req.Request.AddHeader(call.GetArg(0).String(), call.GetArg(1).String())
	return nil
}
func (req *Request) DelHeader(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	req.Request.DelHeader(call.GetArg(0).String())
	return nil

}
func (req *Request) GetHeader(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	return call.Context().NewString(req.Request.GetHeader(call.GetArg(0).String())).Consume()

}
func (req *Request) HeaderValues(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	result := req.Request.HeaderValues(call.GetArg(0).String())
	var output = make([]*v8js.Consumed, len(result))
	for i, v := range result {
		output[i] = call.Context().NewString(v).Consume()
	}
	return call.Context().NewArray(output...).Consume()
}
func (req *Request) HeaderFields(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	result := req.Request.HeaderFields()
	var output = make([]*v8js.Consumed, len(result))
	for i, v := range result {
		output[i] = call.Context().NewString(v).Consume()
	}
	return call.Context().NewArray(output...).Consume()

}

func (req *Request) ResponseStatusCode(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	return call.Context().NewInt32(int32(req.Request.ResponseStatusCode())).Consume()
}
func (req *Request) ResponseBody(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	return call.Context().NewString(string(req.Request.ResponseBody())).Consume()
}
func (req *Request) ResponseHeader(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	return call.Context().NewString(req.Request.ResponseHeader(call.GetArg(0).String())).Consume()

}
func (req *Request) ResponseHeaderValues(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	result := req.Request.ResponseHeaderValues(call.GetArg(0).String())

	var output = make([]*v8js.Consumed, len(result))
	for i, v := range result {
		output[i] = call.Context().NewString(v).Consume()
	}
	return call.Context().NewArray(output...).Consume()
}
func (req *Request) ResponseHeaderFields(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	result := req.Request.ResponseHeaderFields()

	var output = make([]*v8js.Consumed, len(result))
	for i, v := range result {
		output[i] = call.Context().NewString(v).Consume()
	}
	return call.Context().NewArray(output...).Consume()

}
func (req *Request) Execute(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	req.Request.MustExecute()
	return nil
}

type Addon struct {
	Addon   *httpaddon.Addon
	Builder Builder
}

func (a *Addon) ParseURL(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	rawurl := call.Args()[0].String()
	u, err := url.Parse(rawurl)
	if err != nil {
		return nil
	}
	result := call.Context().NewObject()
	result.Set("Host", call.Context().NewString(u.Host).Consume())
	result.Set("Hostname", call.Context().NewString(u.Host).Consume())
	result.Set("Scheme", call.Context().NewString(u.Scheme).Consume())
	result.Set("Path", call.Context().NewString(u.Path).Consume())
	result.Set("Query", call.Context().NewString(u.RawQuery).Consume())
	result.Set("User", call.Context().NewString(u.User.Username()).Consume())
	p, _ := u.User.Password()
	result.Set("Password", call.Context().NewString(p).Consume())
	result.Set("Port", call.Context().NewString(u.Port()).Consume())
	result.Set("Fragment", call.Context().NewString(u.Fragment).Consume())
	return result.Consume()
}
func (a *Addon) NewRequest(call *v8js.FunctionCallbackInfo) *v8js.Consumed {

	method := call.GetArg(0).String()
	url := call.GetArg(1).String()
	req := a.Addon.Create(method, url)
	return a.Builder(call.Context(), &Request{req}).Consume()
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
