package binaryv8

import (
	"github.com/herb-go/herbplugin"
	"github.com/herb-go/plugins/addons/binaryaddon"
	v8js "github.com/jarlyyn/v8js"
)

type Addon struct {
	Addon *binaryaddon.Addon
}

func (a *Addon) Base64Encode(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	data := call.GetArg(0)
	if data != nil {
		return call.Context().NewString(a.Addon.Base64Encode(data.ArrayBufferContent())).Consume()
	}
	return nil
}
func (a *Addon) Base64Decode(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	data := call.GetArg(0)
	if data != nil {
		return call.Context().NewArrayBuffer(a.Addon.Base64Decode(data.String())).Consume()
	}
	return nil
}
func (a *Addon) Md5Sum(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	data := call.GetArg(0)
	if data != nil {
		return call.Context().NewString(a.Addon.Md5Sum(data.ArrayBufferContent())).Consume()
	}
	return nil
}
func (a *Addon) Sha1Sum(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	data := call.GetArg(0)
	if data != nil {
		return call.Context().NewString(a.Addon.Sha1Sum(data.ArrayBufferContent())).Consume()
	}
	return nil
}
func (a *Addon) Sha256Sum(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	data := call.GetArg(0)
	if data != nil {
		return call.Context().NewString(a.Addon.Sha256Sum(data.ArrayBufferContent())).Consume()
	}
	return nil
}
func (a *Addon) Sha512Sum(call *v8js.FunctionCallbackInfo) *v8js.Consumed {
	data := call.GetArg(0)
	if data != nil {
		return call.Context().NewString(a.Addon.Sha512Sum(data.ArrayBufferContent())).Consume()
	}
	return nil
}
func (a *Addon) Convert(r *v8js.Context) *v8js.JsValue {
	obj := r.NewObject()
	obj.SetObjectMethod(r, "Base64Encode", a.Base64Encode)
	obj.SetObjectMethod(r, "Base64Decode", a.Base64Decode)
	obj.SetObjectMethod(r, "Md5Sum", a.Md5Sum)
	obj.SetObjectMethod(r, "Sha1Sum", a.Sha1Sum)
	obj.SetObjectMethod(r, "Sha256Sum", a.Sha256Sum)
	obj.SetObjectMethod(r, "Sha512Sum", a.Sha512Sum)
	return obj
}

func Create(p herbplugin.Plugin) *Addon {
	return &Addon{
		Addon: binaryaddon.Create(p),
	}
}
