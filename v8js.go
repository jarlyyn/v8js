package v8js

import (
	"math/big"
	"runtime"

	"rogchap.com/v8go"
)

func NewContext(opt ...v8go.ContextOption) *Context {
	return &Context{
		Raw: v8go.NewContext(opt...),
	}
}
func WrapContext(c *v8go.Context) *Context {
	return &Context{
		Raw: c,
	}
}

type Context struct {
	Raw *v8go.Context
}

func (c *Context) Global() *JsValue {
	return NewJsValue(c.Raw.Global().Value)
}
func (c *Context) NewStringArray(values ...string) *JsValue {
	args := make([]*JsValue, len(values))
	for i, v := range values {
		args[i] = c.NewValue(v)
	}
	return c.NewArray(args...)
}
func (c *Context) NewArray(values ...*JsValue) *JsValue {
	global := c.Global()
	array := global.Get("Array")
	return array.Call(global, values...)
}
func (c *Context) NewValue(value interface{}) *JsValue {

	v, err := v8go.NewValue(c.Raw.Isolate(), value)
	if err != nil {
		panic(err)
	}
	return NewJsValue(v)
}
func (c *Context) NewObject() *JsValue {
	ot := v8go.NewObjectTemplate(c.Raw.Isolate())
	obj, err := ot.NewInstance(c.Raw)
	result := NewJsValue(obj.Value)
	if err != nil {
		panic(err)
	}
	return result
}
func (c *Context) NewFunctionTemplate(callback FunctionCallback) *FunctionTemplate {
	tmpl := v8go.NewFunctionTemplate(c.Raw.Isolate(), func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		result := callback(newFunctionCallbackInfo(info))
		return result.export()
	})
	return &FunctionTemplate{
		tmpl: tmpl,
	}
}
func (c *Context) RunScript(script string, name string) *JsValue {
	result, err := c.Raw.RunScript(script, name)
	if err != nil {
		panic(err)
	}
	return NewJsValue(result)
}
func (c *Context) CloseAndDispose() {
	c.Raw.Close()
	c.Raw.Isolate().Dispose()
}
func NewJsValue(raw *v8go.Value) *JsValue {
	v := &JsValue{
		Raw:      raw,
		Exported: false,
	}
	runtime.SetFinalizer(v, v.Release)
	return v
}

type JsValue struct {
	Raw      *v8go.Value
	Exported bool
}

func mustAsObject(v *v8go.Value) *v8go.Object {
	o, err := v.AsObject()
	if err != nil {
		panic(err)
	}
	return o
}
func (v *JsValue) export() *v8go.Value {
	v.Exported = true
	return v.Raw
}
func (v *JsValue) Call(recvr *JsValue, args ...*JsValue) *JsValue {
	fn, err := v.Raw.AsFunction()
	if err != nil {
		panic(err)
	}
	fnargs := make([]v8go.Valuer, len(args))
	for i, v := range args {
		fnargs[i] = v.export()
	}
	val, err := fn.Call(recvr.export(), fnargs...)
	result := NewJsValue(val)
	if err != nil {
		panic(err)
	}
	return result
}

func (v *JsValue) Release() {
	if !v.Exported {
		v.Exported = true
		if v.Raw != nil {
			v.Raw.Release()
		}
	}
}

func (v *JsValue) String() string {
	return v.Raw.String()
}

func (v *JsValue) BigInt() *big.Int {
	return v.Raw.BigInt()
}
func (v *JsValue) Boolean() bool {
	return v.Raw.Boolean()
}
func (v *JsValue) Int32() int32 {
	return v.Raw.Int32()
}
func (v *JsValue) Integer() int64 {
	return v.Raw.Integer()
}
func (v *JsValue) Number() float64 {
	return v.Raw.Number()
}
func (v *JsValue) Uint32() uint32 {
	return v.Raw.Uint32()
}
func (v *JsValue) SameValue(other *JsValue) bool {
	return v.Raw.SameValue(other.Raw)
}
func (v *JsValue) IsUndefined() bool {
	return v.Raw.IsUndefined()
}
func (v *JsValue) IsNull() bool {
	return v.Raw.IsNull()
}
func (v *JsValue) IsNullOrUndefined() bool {
	return v.Raw.IsNullOrUndefined()
}

func (v *JsValue) IsTrue() bool {
	return v.Raw.IsTrue()
}
func (v *JsValue) IsFalse() bool {
	return v.Raw.IsFalse()
}

func (v *JsValue) IsFunction() bool {
	return v.Raw.IsFunction()
}

func (v *JsValue) IsObject() bool {
	return v.Raw.IsObject()
}

func (v *JsValue) IsBigInt() bool {
	return v.Raw.IsBigInt()
}
func (v *JsValue) IsBoolean() bool {
	return v.Raw.IsBoolean()
}

func (v *JsValue) IsNumber() bool {
	return v.Raw.IsNumber()
}
func (v *JsValue) IsInt32() bool {
	return v.Raw.IsInt32()
}
func (v *JsValue) IsUint32() bool {
	return v.Raw.IsUint32()
}
func (v *JsValue) IsDate() bool {
	return v.Raw.IsDate()
}
func (v *JsValue) IsNativeError() bool {
	return v.Raw.IsNativeError()
}
func (v *JsValue) IsRegExp() bool {
	return v.Raw.IsRegExp()
}
func (v *JsValue) IsMap() bool {
	return v.Raw.IsMap()
}
func (v *JsValue) IsSet() bool {
	return v.Raw.IsSet()
}
func (v *JsValue) IsArray() bool {
	return v.Raw.IsArray()
}

func (v *JsValue) MustMarshalJSON() []byte {
	data, err := v.Raw.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return data
}

func (v *JsValue) MethodCall(methodName string, args ...*JsValue) *JsValue {
	fn := v.Get(methodName) // ensure method exists
	return fn.Call(v, args...)
}

func (v *JsValue) Get(key string) *JsValue {
	val, err := mustAsObject(v.Raw).Get(key)
	if err != nil {
		panic(err)
	}
	return NewJsValue(val)
}
func (v *JsValue) GetIdx(idx uint32) *JsValue {
	val, err := mustAsObject(v.Raw).GetIdx(idx)
	if err != nil {
		panic(err)
	}
	return NewJsValue(val)
}

func (v *JsValue) Set(key string, val *JsValue) {
	err := mustAsObject(v.Raw).Set(key, val.export())
	if err != nil {
		panic(err)
	}
}

func (v *JsValue) SetIdx(idx uint32, val *JsValue) {
	err := mustAsObject(v.Raw).SetIdx(idx, val.export())
	if err != nil {
		panic(err)
	}
}
func (v *JsValue) Has(key string) bool {
	return mustAsObject(v.Raw).Has(key)
}
func (v *JsValue) HasIdx(idx uint32) bool {
	return mustAsObject(v.Raw).HasIdx(idx)
}

func (v *JsValue) Delete(key string) bool {
	return mustAsObject(v.Raw).Delete(key)
}

func (v *JsValue) DeleteIdx(idx uint32) bool {
	return mustAsObject(v.Raw).DeleteIdx(idx)
}

type FunctionCallback func(info *FunctionCallbackInfo) *JsValue

type FunctionCallbackInfo struct {
	ctx  *Context
	args []*JsValue
	this *JsValue
}

func newFunctionCallbackInfo(raw *v8go.FunctionCallbackInfo) *FunctionCallbackInfo {
	rawargs := raw.Args()
	args := make([]*JsValue, len(rawargs))
	for k, v := range rawargs {
		args[k] = NewJsValue(v)
	}
	return &FunctionCallbackInfo{
		ctx:  WrapContext(raw.Context()),
		args: args,
		this: NewJsValue(raw.This().Value),
	}
}
func (i *FunctionCallbackInfo) Context() *Context {
	return i.ctx
}

// This returns the receiver object "this".
func (i *FunctionCallbackInfo) This() *JsValue {
	return i.this
}

// Args returns a slice of the value arguments that are passed to the JS function.
func (i *FunctionCallbackInfo) Args() []*JsValue {
	return i.args
}

type FunctionTemplate struct {
	tmpl *v8go.FunctionTemplate
}

func (t *FunctionTemplate) GetFunction(ctx *Context) *JsValue {
	fn := t.tmpl.GetFunction(ctx.Raw)
	return NewJsValue(fn.Value)
}
