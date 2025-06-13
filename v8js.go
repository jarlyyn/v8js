package v8js

import (
	"math/big"
	"runtime"

	"rogchap.com/v8go"
)

func NewContext(opt ...v8go.ContextOption) *Context {
	c := &Context{Raw: v8go.NewContext(opt...)}
	c.objectTemplate = v8go.NewObjectTemplate(c.Raw.Isolate())
	return c

}

type Context struct {
	Raw            *v8go.Context
	objectTemplate *v8go.ObjectTemplate
}

func (c *Context) Global() *JsValue {
	return PromiseReleaseInFuture(c.Raw.Global().Value)
}

func (c *Context) newValue(v interface{}) *JsValue {
	val, err := v8go.NewValue(c.Raw.Isolate(), v)
	if err != nil {
		panic(err)
	}
	return PromiseReleaseInFuture(val)
}
func (c *Context) NewString(val string) *JsValue {
	return c.newValue(val)
}
func (c *Context) NewInt32(val int32) *JsValue {
	return c.newValue(val)
}
func (c *Context) NewInt64(val int64) *JsValue {
	return c.newValue(val)
}
func (c *Context) NewBoolean(val bool) *JsValue {
	return c.newValue(val)
}
func (c *Context) NewBigInt(val *big.Int) *JsValue {
	return c.newValue(val)
}
func (c *Context) NewNumber(val float64) *JsValue {
	return c.newValue(val)
}
func (c *Context) NewFunction(callback FunctionCallback) *JsValue {
	tmpl := c.NewFunctionTemplate(callback)
	fn := tmpl.GetFunction(c)
	return fn
}

func (c *Context) NewStringArray(values ...string) *JsValue {
	args := make([]*JsValue, len(values))
	for i, v := range values {
		args[i] = c.NewString(v)
	}
	return c.NewArray(args...)
}
func (c *Context) NewArray(values ...*JsValue) *JsValue {
	global := c.Global()
	array := global.Get("Array")
	return array.Call(global, values...)
}
func (c *Context) NewObject() *JsValue {
	obj, err := c.objectTemplate.NewInstance(c.Raw)
	if err != nil {
		panic(err)
	}
	result := PromiseReleaseInFuture(obj.Value)
	return result
}
func (c *Context) NewFunctionTemplate(callback FunctionCallback) *FunctionTemplate {
	tmpl := v8go.NewFunctionTemplate(c.Raw.Isolate(), func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		result := callback(newFunctionCallbackInfo(c, info))
		defer runtime.KeepAlive(result)
		if result == nil {
			return nil
		}
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
	return PromiseReleaseInFuture(result)
}
func (c *Context) CloseAndDispose() {
	c.Raw.Close()
	c.Raw.Isolate().Dispose()
}
func (c *Context) NullValue() *JsValue {
	val := PromiseReleaseInFuture(v8go.Null(c.Raw.Isolate()))
	val.noRelease = true
	return val
}

type JsValue struct {
	raw       *v8go.Value
	noRelease bool
}

func mustAsObject(v *v8go.Value) *v8go.Object {
	o, err := v.AsObject()
	if err != nil {
		panic(err)
	}
	return o
}
func (v *JsValue) export() *v8go.Value {
	if v == nil {
		return nil
	}
	return v.raw
}
func (v *JsValue) Call(recvr *JsValue, args ...*JsValue) *JsValue {
	defer runtime.KeepAlive(v)
	defer runtime.KeepAlive(recvr)

	fn, err := v.raw.AsFunction()

	if err != nil {
		panic(err)
	}
	fnargs := make([]v8go.Valuer, len(args))
	for i, v := range args {
		fnargs[i] = v.export()
		defer runtime.KeepAlive(v)
	}
	val, err := fn.Call(recvr.export(), fnargs...)
	if err != nil {
		panic(err)
	}
	result := PromiseReleaseInFuture(val)

	return result
}

func (v *JsValue) Release() {
	if v.raw != nil {
		v.raw = nil
		if !v.noRelease {
			v.raw.Release()
		}
	}
}

func (v *JsValue) Array() []*JsValue {
	result := []*JsValue{}
	length := v.Get("length")
	if length.IsNullOrUndefined() {
		return result
	}
	ln := int(length.Integer())
	for i := 0; i < ln; i++ {
		item := v.GetIdx(uint32(i))
		if item.IsNullOrUndefined() {
			continue
		}
		result = append(result, item)
	}
	return result
}
func (v *JsValue) StringArrry() []string {
	arr := v.Array()
	result := make([]string, len(arr))
	for i, item := range arr {
		result[i] = item.String()
	}
	return result
}
func (v *JsValue) String() string {
	return v.raw.String()
}

func (v *JsValue) BigInt() *big.Int {
	return v.raw.BigInt()
}
func (v *JsValue) Boolean() bool {
	return v.raw.Boolean()
}
func (v *JsValue) Int32() int32 {
	return v.raw.Int32()
}
func (v *JsValue) Integer() int64 {
	return v.raw.Integer()
}
func (v *JsValue) Number() float64 {
	return v.raw.Number()
}
func (v *JsValue) Uint32() uint32 {
	return v.raw.Uint32()
}
func (v *JsValue) SameValue(other *JsValue) bool {
	return v.raw.SameValue(other.raw)
}
func (v *JsValue) IsUndefined() bool {
	return v.raw.IsUndefined()
}
func (v *JsValue) IsNull() bool {
	return v.raw.IsNull()
}
func (v *JsValue) IsNullOrUndefined() bool {
	return v.raw.IsNullOrUndefined()
}

func (v *JsValue) IsTrue() bool {
	return v.raw.IsTrue()
}
func (v *JsValue) IsFalse() bool {
	return v.raw.IsFalse()
}

func (v *JsValue) IsFunction() bool {
	return v.raw.IsFunction()
}

func (v *JsValue) IsObject() bool {
	return v.raw.IsObject()
}

func (v *JsValue) IsBigInt() bool {
	return v.raw.IsBigInt()
}
func (v *JsValue) IsBoolean() bool {
	return v.raw.IsBoolean()
}

func (v *JsValue) IsNumber() bool {
	return v.raw.IsNumber()
}
func (v *JsValue) IsInt32() bool {
	return v.raw.IsInt32()
}
func (v *JsValue) IsUint32() bool {
	return v.raw.IsUint32()
}
func (v *JsValue) IsDate() bool {
	return v.raw.IsDate()
}
func (v *JsValue) IsNativeError() bool {
	return v.raw.IsNativeError()
}
func (v *JsValue) IsRegExp() bool {
	return v.raw.IsRegExp()
}
func (v *JsValue) IsMap() bool {
	return v.raw.IsMap()
}
func (v *JsValue) IsSet() bool {
	return v.raw.IsSet()
}
func (v *JsValue) IsArray() bool {
	return v.raw.IsArray()
}

func (v *JsValue) MustMarshalJSON() []byte {
	data, err := v.raw.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return data
}

func (v *JsValue) MethodCall(methodName string, args ...*JsValue) *JsValue {
	fn := v.Get(methodName) // ensure method exists
	return fn.Call(v, args...)
}
func (v *JsValue) SetObjectMethod(ctx *Context, name string, fn FunctionCallback) {
	f := ctx.NewFunctionTemplate(fn).GetFunction(ctx)
	v.Set(name, f)
}
func (v *JsValue) Get(key string) *JsValue {
	val, err := mustAsObject(v.raw).Get(key)
	if err != nil {
		panic(err)
	}
	return PromiseReleaseInFuture(val)
}
func (v *JsValue) GetIdx(idx uint32) *JsValue {
	val, err := mustAsObject(v.raw).GetIdx(idx)
	if err != nil {
		panic(err)
	}
	return PromiseReleaseInFuture(val)
}

func (v *JsValue) Set(key string, val *JsValue) {
	defer runtime.KeepAlive(val)
	err := mustAsObject(v.raw).Set(key, val.export())
	if err != nil {
		panic(err)
	}
}

func (v *JsValue) SetIdx(idx uint32, val *JsValue) {
	defer runtime.KeepAlive(val)
	err := mustAsObject(v.raw).SetIdx(idx, val.export())
	if err != nil {
		panic(err)
	}
}
func (v *JsValue) Has(key string) bool {
	return mustAsObject(v.raw).Has(key)
}
func (v *JsValue) HasIdx(idx uint32) bool {
	return mustAsObject(v.raw).HasIdx(idx)
}

func (v *JsValue) Delete(key string) bool {
	return mustAsObject(v.raw).Delete(key)
}

func (v *JsValue) DeleteIdx(idx uint32) bool {
	return mustAsObject(v.raw).DeleteIdx(idx)
}

type FunctionCallback func(info *FunctionCallbackInfo) *JsValue

type FunctionCallbackInfo struct {
	ctx  *Context
	args []*JsValue
	this *JsValue
}

func newFunctionCallbackInfo(ctx *Context, raw *v8go.FunctionCallbackInfo) *FunctionCallbackInfo {
	rawargs := raw.Args()
	args := make([]*JsValue, len(rawargs))
	for k, v := range rawargs {
		args[k] = PromiseReleaseInFuture(v)
	}
	return &FunctionCallbackInfo{
		ctx:  ctx,
		args: args,
		this: PromiseReleaseInFuture(raw.This().Value),
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
func (i *FunctionCallbackInfo) GetArg(idx int) *JsValue {
	if idx < 0 || idx >= len(i.args) {
		return i.ctx.NullValue()
	}
	return i.args[idx]
}

type FunctionTemplate struct {
	tmpl *v8go.FunctionTemplate
}

func (t *FunctionTemplate) GetFunction(ctx *Context) *JsValue {
	fn := t.tmpl.GetFunction(ctx.Raw)
	return PromiseReleaseInFuture(fn.Value)
}

func ReleaseJsValueAsRawValue(v *JsValue) *v8go.Value {
	val := v.raw
	v.raw = nil
	return val
}

func PromiseReleaseInFuture(v *v8go.Value) *JsValue {
	val := &JsValue{
		raw: v,
	}
	runtime.SetFinalizer(val, func(jv *JsValue) { jv.Release() })
	return val
}
