package v8js

import (
	"math/big"
	"runtime"
	"sync"

	"github.com/herb-go/v8go"
)

func NewContext(opt ...v8go.ContextOption) *Context {
	c := &Context{Raw: v8go.NewContext(opt...)}
	c.objectTemplate = v8go.NewObjectTemplate(c.Raw.Isolate())
	c.nullvalue = c.WrapWithoutReleaser(v8go.Null(c.Raw.Isolate()))
	return c
}

type Context struct {
	Locker         sync.Mutex
	Raw            *v8go.Context
	objectTemplate *v8go.ObjectTemplate
	nullvalue      *JsValue
}

func (c *Context) Close() {
	c.Locker.Lock()
	defer c.Locker.Unlock()
	if c.Raw == nil {
		return
	}
	c.Raw.Close()
	c.Raw.Isolate().Dispose()
	c.Raw = nil
	c.objectTemplate = nil
	c.nullvalue = nil
}
func (c *Context) WrapWithoutReleaser(v *v8go.Value) *JsValue {
	val := &JsValue{
		raw:       v,
		ctx:       c,
		noRelease: true,
	}
	return val
}
func (c *Context) PromiseReleaseInFuture(v *v8go.Value) *JsValue {
	val := &JsValue{
		raw: v,
		ctx: c,
	}
	runtime.SetFinalizer(val, (*JsValue).Release)
	return val
}
func (c *Context) Global() *JsValue {
	result := c.PromiseReleaseInFuture(c.Raw.Global().Value)
	return result
}

func (c *Context) newValue(v interface{}) *JsValue {
	val, err := v8go.NewValue(c.Raw.Isolate(), v)
	if err != nil {
		panic(err)
	}
	return c.PromiseReleaseInFuture(val)
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
	result := c.PromiseReleaseInFuture(obj.Value) //?
	return result
}
func (c *Context) NewFunctionTemplate(callback FunctionCallback) *FunctionTemplate {
	return newFunctionTemplate(c, callback)
}
func (c *Context) RunScript(script string, name string) *JsValue {
	result, err := c.Raw.RunScript(script, name)
	if err != nil {
		panic(err)
	}
	return c.PromiseReleaseInFuture(result) //报错严重
}
func (c *Context) NullValue() *JsValue {
	return c.nullvalue
}

type JsValue struct {
	raw       *v8go.Value
	ctx       *Context
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
	result := v.raw
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) Call(recvr *JsValue, args ...*JsValue) *JsValue {
	fn, err := v.export().AsFunction()

	if err != nil {
		panic(err)
	}
	fnargs := make([]v8go.Valuer, len(args))
	for i, val := range args {
		fnargs[i] = val.export()
	}
	val, err := fn.Call(recvr.export(), fnargs...)
	if err != nil {
		panic(err)
	}
	result := v.ctx.PromiseReleaseInFuture(val) //必然崩溃点
	runtime.KeepAlive(v)
	runtime.KeepAlive(recvr)
	runtime.KeepAlive(args)
	return result
}

func (v *JsValue) Release() {
	runtime.SetFinalizer(v, nil) // remove finalizer to prevent double release
	if v.raw != nil {
		ptr := v.raw
		v.raw = nil
		if !v.noRelease {
			v.ctx.Locker.Lock()
			defer v.ctx.Locker.Unlock()
			if v.ctx.Raw == nil {
				return
			}
			ptr.Release()
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
	result := v.export().String()
	runtime.KeepAlive(v)
	return result
}

func (v *JsValue) BigInt() *big.Int {
	result := v.export().BigInt()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) Boolean() bool {
	result := v.export().Boolean()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) Int32() int32 {
	result := v.export().Int32()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) Integer() int64 {
	result := v.export().Integer()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) Number() float64 {
	result := v.export().Number()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) Uint32() uint32 {
	result := v.export().Uint32()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) SameValue(other *JsValue) bool {
	result := v.export().SameValue(other.raw)
	runtime.KeepAlive(other)
	return result
}
func (v *JsValue) IsUndefined() bool {
	result := v.export().IsUndefined()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) IsNull() bool {
	result := v.export().IsNull()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) IsNullOrUndefined() bool {
	result := v.export().IsNullOrUndefined()
	runtime.KeepAlive(v)
	return result
}

func (v *JsValue) IsTrue() bool {
	result := v.export().IsTrue()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) IsFalse() bool {
	result := v.export().IsFalse()
	runtime.KeepAlive(v)
	return result
}

func (v *JsValue) IsFunction() bool {
	result := v.export().IsFunction()
	runtime.KeepAlive(v)
	return result
}

func (v *JsValue) IsObject() bool {
	result := v.export().IsObject()
	runtime.KeepAlive(v)
	return result
}

func (v *JsValue) IsBigInt() bool {
	result := v.export().IsBigInt()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) IsBoolean() bool {
	result := v.export().IsBoolean()
	runtime.KeepAlive(v)
	return result
}

func (v *JsValue) IsNumber() bool {
	result := v.export().IsNumber()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) IsInt32() bool {
	result := v.export().IsInt32()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) IsUint32() bool {
	result := v.export().IsUint32()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) IsDate() bool {
	result := v.export().IsDate()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) IsNativeError() bool {
	result := v.export().IsNativeError()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) IsRegExp() bool {
	result := v.export().IsRegExp()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) IsMap() bool {
	result := v.export().IsMap()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) IsSet() bool {
	result := v.export().IsSet()
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) IsArray() bool {
	result := v.export().IsArray()
	runtime.KeepAlive(v)
	return result
}

func (v *JsValue) MustMarshalJSON() []byte {
	data, err := v.export().MarshalJSON()
	if err != nil {
		panic(err)
	}
	runtime.KeepAlive(v)
	return data
}

func (v *JsValue) MethodCall(methodName string, args ...*JsValue) *JsValue {
	fn := v.Get(methodName) // ensure method exists
	result := fn.Call(v, args...)
	runtime.KeepAlive(fn)
	runtime.KeepAlive(args)
	return result
}
func (v *JsValue) SetObjectMethod(ctx *Context, name string, fn FunctionCallback) {
	f := ctx.NewFunctionTemplate(fn).GetFunction(ctx)
	v.Set(name, f)
	runtime.KeepAlive(f)
}
func (v *JsValue) Get(key string) *JsValue {
	val, err := mustAsObject(v.export()).Get(key)
	if err != nil {
		panic(err)
	}
	result := v.ctx.PromiseReleaseInFuture(val)
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) GetIdx(idx uint32) *JsValue {
	val, err := mustAsObject(v.export()).GetIdx(idx)
	if err != nil {
		panic(err)
	}
	result := v.ctx.PromiseReleaseInFuture(val)
	runtime.KeepAlive(v)
	return result

}

func (v *JsValue) Set(key string, val *JsValue) {
	err := mustAsObject(v.export()).Set(key, val.export())
	if err != nil {
		panic(err)
	}
	runtime.KeepAlive(val)
	runtime.KeepAlive(v)
}

func (v *JsValue) SetIdx(idx uint32, val *JsValue) {
	err := mustAsObject(v.export()).SetIdx(idx, val.export())
	if err != nil {
		panic(err)
	}
	runtime.KeepAlive(val)
	runtime.KeepAlive(v)
}
func (v *JsValue) Has(key string) bool {
	result := mustAsObject(v.export()).Has(key)
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) HasIdx(idx uint32) bool {
	result := mustAsObject(v.export()).HasIdx(idx)
	runtime.KeepAlive(v)
	return result
}

func (v *JsValue) Delete(key string) bool {
	result := mustAsObject(v.export()).Delete(key)
	runtime.KeepAlive(v)
	return result
}

func (v *JsValue) DeleteIdx(idx uint32) bool {
	result := mustAsObject(v.export()).DeleteIdx(idx)
	runtime.KeepAlive(v)
	return result
}

type FunctionCallback func(info *FunctionCallbackInfo) *JsValue

func (f FunctionCallback) convert(c *Context) v8go.FunctionCallback {
	return func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		rawargs := info.Args()
		args := make([]*JsValue, len(rawargs))
		for k, v := range rawargs {
			args[k] = c.WrapWithoutReleaser(v)
		}
		this := c.WrapWithoutReleaser(info.This().Value)
		fi := &FunctionCallbackInfo{
			ctx:  c,
			args: args,
			this: this,
		}
		result := f(fi)
		if result == nil {
			return nil
		}
		result.noRelease = true
		output := result.export()
		runtime.KeepAlive(result)
		runtime.KeepAlive(this)
		runtime.KeepAlive(result)
		runtime.KeepAlive(args)
		runtime.KeepAlive(info)
		runtime.KeepAlive(fi)
		return output
	}
}

type FunctionCallbackInfo struct {
	ctx  *Context
	args []*JsValue
	this *JsValue
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
	return ctx.PromiseReleaseInFuture(fn.Value)
}
func newFunctionTemplate(c *Context, callback FunctionCallback) *FunctionTemplate {
	tmpl := v8go.NewFunctionTemplate(c.Raw.Isolate(), callback.convert(c))
	return &FunctionTemplate{
		tmpl: tmpl,
	}

}
func ExportRawValue(v *JsValue, noRlease bool) *v8go.Value {
	val := v.raw
	if noRlease {
		v.noRelease = true
	}
	runtime.KeepAlive(v)
	return val
}
