package httpv8

import (
	"net/url"
	"strconv"
	"sync"
	"unsafe"

	"github.com/herb-go/herbplugin"
	"github.com/herb-go/plugins/addons/httpaddon"
	v8js "github.com/jarlyyn/v8js"
)

type Builder func(r *v8js.Context, a *Addon, req *Request) *v8js.JsValue

var DefaultBuilder = func(r *v8js.Context, a *Addon, req *Request) *v8js.JsValue {
	obj := r.NewObject()
	obj.Set("GetID", a.Functions["GetID"].Consume())
	obj.Set("GetURL", a.Functions["GetURL"].Consume())
	obj.Set("SetURL", a.Functions["SetURL"].Consume())
	obj.Set("GetProxy", a.Functions["GetProxy"].Consume())
	obj.Set("SetProxy", a.Functions["SetProxy"].Consume())
	obj.Set("GetMethod", a.Functions["GetMethod"].Consume())
	obj.Set("SetMethod", a.Functions["SetMethod"].Consume())
	obj.Set("GetBody", a.Functions["GetBody"].Consume())
	obj.Set("GetBodyArrayBuffer", a.Functions["GetBodyArrayBuffer"].Consume())
	obj.Set("SetBody", a.Functions["SetBody"].Consume())
	obj.Set("FinishedAt", a.Functions["FinishedAt"].Consume())
	obj.Set("ExecuteStatus", a.Functions["ExecuteStatus"].Consume())
	obj.Set("ResetHeader", a.Functions["ResetHeader"].Consume())
	obj.Set("SetHeader", a.Functions["SetHeader"].Consume())
	obj.Set("AddHeader", a.Functions["AddHeader"].Consume())
	obj.Set("DelHeader", a.Functions["DelHeader"].Consume())
	obj.Set("GetHeader", a.Functions["GetHeader"].Consume())
	obj.Set("HeaderValues", a.Functions["HeaderValues"].Consume())
	obj.Set("HeaderFields", a.Functions["HeaderFields"].Consume())
	obj.Set("ResponseStatusCode", a.Functions["ResponseStatusCode"].Consume())
	obj.Set("ResponseBody", a.Functions["ResponseBody"].Consume())
	obj.Set("ResponseBodyArrayBuffer", a.Functions["ResponseBodyArrayBuffer"].Consume())
	obj.Set("ResponseHeader", a.Functions["ResponseHeader"].Consume())
	obj.Set("ResponseHeaderValues", a.Functions["ResponseHeaderValues"].Consume())
	obj.Set("ResponseHeaderFields", a.Functions["ResponseHeaderFields"].Consume())
	obj.Set("Execute", a.Functions["Execute"].Consume())
	return obj
}

type Request struct {
	RID     string
	Request *httpaddon.Request
}

func RequestGetID(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		return call.Context().NewString(req.Request.GetID()).Consume()
	}
}

func RequestGetURL(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		return call.Context().NewString(req.Request.GetURL()).Consume()
	}
}
func RequestSetURL(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		req.Request.SetURL(call.GetArg(0).String())
		return nil
	}
}
func RequestGetProxy(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		return call.Context().NewString(req.Request.GetProxy()).Consume()
	}
}
func RequestSetProxy(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		req.Request.SetProxy(call.GetArg(0).String())
		return nil
	}
}

func RequestGetMethod(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		return call.Context().NewString(req.Request.GetMethod()).Consume()
	}
}
func RequestSetMethod(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		req.Request.SetMethod(call.GetArg(0).String())
		return nil
	}
}
func RequestGetBody(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		return call.Context().NewString(string(req.Request.GetBody())).Consume()
	}
}
func RequestGetBodyArrayBuffer(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		return call.Context().NewArrayBuffer(req.Request.GetBody()).Consume()
	}
}

func RequestSetBody(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		req.Request.SetBody([]byte(call.GetArg(0).String()))
		return nil
	}
}

func RequestFinishedAt(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		return call.Context().NewInt64(req.Request.FinishedAt()).Consume()
	}
}
func RequestExecuteStatus(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		return call.Context().NewInt32(int32(req.Request.ExecuteStauts())).Consume()
	}
}
func RequestResetHeader(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		req.Request.ResetHeader()
		return nil
	}
}
func RequestSetHeader(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		req.Request.SetHeader(call.GetArg(0).String(), call.GetArg(1).String())
		return nil
	}
}
func RequestAddHeader(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		req.Request.AddHeader(call.GetArg(0).String(), call.GetArg(1).String())
		return nil
	}
}
func RequestDelHeader(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		req.Request.DelHeader(call.GetArg(0).String())
		return nil
	}

}
func RequestGetHeader(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		return call.Context().NewString(req.Request.GetHeader(call.GetArg(0).String())).Consume()
	}
}
func RequestHeaderValues(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		result := req.Request.HeaderValues(call.GetArg(0).String())
		var output = make([]*v8js.Consumed, len(result))
		for i, v := range result {
			output[i] = call.Context().NewString(v).Consume()
		}
		return call.Context().NewArray(output...).Consume()
	}
}
func RequestHeaderFields(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		result := req.Request.HeaderFields()
		var output = make([]*v8js.Consumed, len(result))
		for i, v := range result {
			output[i] = call.Context().NewString(v).Consume()
		}
		return call.Context().NewArray(output...).Consume()
	}
}

func RequestResponseStatusCode(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		return call.Context().NewInt32(int32(req.Request.ResponseStatusCode())).Consume()
	}
}
func RequestResponseBody(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		return call.Context().NewString(string(req.Request.ResponseBody())).Consume()
	}
}
func RequestResponseBodyArrayBuffer(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		return call.Context().NewArrayBuffer(req.Request.ResponseBody()).Consume()
	}
}
func RequestResponseHeader(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		return call.Context().NewString(req.Request.ResponseHeader(call.GetArg(0).String())).Consume()
	}
}
func RequestResponseHeaderValues(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		result := req.Request.ResponseHeaderValues(call.GetArg(0).String())

		var output = make([]*v8js.Consumed, len(result))
		for i, v := range result {
			output[i] = call.Context().NewString(v).Consume()
		}
		return call.Context().NewArray(output...).Consume()
	}
}
func RequestResponseHeaderFields(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		result := req.Request.ResponseHeaderFields()

		var output = make([]*v8js.Consumed, len(result))
		for i, v := range result {
			output[i] = call.Context().NewString(v).Consume()
		}
		return call.Context().NewArray(output...).Consume()
	}
}
func RequestExecute(a *Addon) func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	return func(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
		req := a.LoadReq(call.This().Get("id").String())
		req.Request.MustExecute()
		return nil
	}
}

type Addon struct {
	Addon     *httpaddon.Addon
	Builder   Builder
	Functions map[string]*v8js.Reusable
	reqs      sync.Map
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
	rid := strconv.FormatInt(int64(uintptr(unsafe.Pointer(req))), 16)
	ar := &Request{
		RID:     rid,
		Request: req,
	}
	a.reqs.Store(rid, ar)
	fr := call.This().Get("FinalizationRegistry")
	obj := a.Builder(call.Context(), a, ar)
	obj.Set("id", call.Context().NewString(rid).Consume())
	oc := obj.ConsumeReuseble()

	fr.Get("register").Call(fr, oc.Consume(), a.Register(call.Context(), call.This(), req.ID).Consume())
	return oc.FinalConsume()
}
func (a *Addon) Register(r *v8js.Context, addonobj *v8js.Consumed, id string) *v8js.JsValue {
	obj := r.NewObject()
	obj.Set("unload", addonobj.Get("unload").Consume())
	obj.Set("id", r.NewString(id).Consume())
	return obj
}
func (a *Addon) unload(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	a.reqs.Delete(call.GetArg(0).String())
	return nil
}
func (a *Addon) size(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	count := 0
	a.reqs.Range(func(key, value interface{}) bool {
		count++
		return true // continue iteration
	})

	return call.Context().NewInt32(int32(count)).Consume()
}
func (a *Addon) LoadReq(id string) *Request {
	v, ok := a.reqs.Load(id)
	if !ok {
		panic("v8 http request id " + id + " not found")
	}
	return v.(*Request)
}
func (a *Addon) Convert(r *v8js.Context) *v8js.JsValue {
	obj := r.NewObject()
	obj.SetObjectMethod(r, "New", a.NewRequest)
	obj.SetObjectMethod(r, "ParseURL", a.ParseURL)
	fr := r.RunScript(" new FinalizationRegistry((reg) => {reg.unload(reg.id)})", "FinalizationRegistry")
	obj.Set("FinalizationRegistry", fr.Consume())
	obj.Set("unload", r.NewFunction(a.unload).Consume())
	obj.SetObjectMethod(r, "Size", a.size)
	a.Functions["GetID"] = r.NewFunction(RequestGetID(a)).ConsumeReuseble()
	a.Functions["GetURL"] = r.NewFunction(RequestGetURL(a)).ConsumeReuseble()
	a.Functions["SetURL"] = r.NewFunction(RequestSetURL(a)).ConsumeReuseble()
	a.Functions["GetProxy"] = r.NewFunction(RequestGetProxy(a)).ConsumeReuseble()
	a.Functions["SetProxy"] = r.NewFunction(RequestSetProxy(a)).ConsumeReuseble()
	a.Functions["GetMethod"] = r.NewFunction(RequestGetMethod(a)).ConsumeReuseble()
	a.Functions["SetMethod"] = r.NewFunction(RequestSetMethod(a)).ConsumeReuseble()
	a.Functions["GetBody"] = r.NewFunction(RequestGetBody(a)).ConsumeReuseble()
	a.Functions["GetBodyArrayBuffer"] = r.NewFunction(RequestGetBodyArrayBuffer(a)).ConsumeReuseble()
	a.Functions["SetBody"] = r.NewFunction(RequestSetBody(a)).ConsumeReuseble()
	a.Functions["FinishedAt"] = r.NewFunction(RequestFinishedAt(a)).ConsumeReuseble()
	a.Functions["ExecuteStatus"] = r.NewFunction(RequestExecuteStatus(a)).ConsumeReuseble()
	a.Functions["ResetHeader"] = r.NewFunction(RequestResetHeader(a)).ConsumeReuseble()
	a.Functions["SetHeader"] = r.NewFunction(RequestSetHeader(a)).ConsumeReuseble()
	a.Functions["AddHeader"] = r.NewFunction(RequestAddHeader(a)).ConsumeReuseble()
	a.Functions["DelHeader"] = r.NewFunction(RequestDelHeader(a)).ConsumeReuseble()
	a.Functions["GetHeader"] = r.NewFunction(RequestGetHeader(a)).ConsumeReuseble()
	a.Functions["HeaderValues"] = r.NewFunction(RequestHeaderValues(a)).ConsumeReuseble()
	a.Functions["HeaderFields"] = r.NewFunction(RequestHeaderFields(a)).ConsumeReuseble()
	a.Functions["ResponseStatusCode"] = r.NewFunction(RequestResponseStatusCode(a)).ConsumeReuseble()
	a.Functions["ResponseBody"] = r.NewFunction(RequestResponseBody(a)).ConsumeReuseble()
	a.Functions["ResponseBodyArrayBuffer"] = r.NewFunction(RequestResponseBodyArrayBuffer(a)).ConsumeReuseble()
	a.Functions["ResponseHeader"] = r.NewFunction(RequestResponseHeader(a)).ConsumeReuseble()
	a.Functions["ResponseHeaderValues"] = r.NewFunction(RequestResponseHeaderValues(a)).ConsumeReuseble()
	a.Functions["ResponseHeaderFields"] = r.NewFunction(RequestResponseHeaderFields(a)).ConsumeReuseble()
	a.Functions["Execute"] = r.NewFunction(RequestExecute(a)).ConsumeReuseble()

	return obj
}
func Create(p herbplugin.Plugin) *Addon {
	return &Addon{
		Addon:     httpaddon.Create(p),
		Functions: make(map[string]*v8js.Reusable),
		Builder:   DefaultBuilder,
	}
}
