package v8js

import (
	"fmt"
	"math/big"
	"runtime"
	"sync"

	"github.com/herb-go/v8go"
)

func NewContext(opt ...v8go.ContextOption) *Context {
	c := &Context{Raw: v8go.NewContext(opt...)}
	c.objectTemplate = v8go.NewObjectTemplate(c.Raw.Isolate())
	c.nullvalue = c.Wrap(v8go.Null(c.Raw.Isolate()))
	return c
}

type Context struct {
	locker         sync.RWMutex
	objectTemplate *v8go.ObjectTemplate
	Raw            *v8go.Context
	nullvalue      *JsValue
}

func (c *Context) Close() {
	c.locker.Lock()
	defer c.locker.Unlock()
	if c.Raw == nil {
		return
	}
	ctx := c.Raw
	c.Raw = nil
	c.nullvalue = nil
	c.objectTemplate = nil
	ctx.Close()
	ctx.Isolate().Dispose()
	runtime.GC()
}
func (c *Context) Wrap(v *v8go.Value) *JsValue {
	val := &JsValue{
		raw: v,
		ctx: c,
	}
	return val
}

func (c *Context) Global() *JsValue {
	result := c.Wrap(c.Raw.Global().Value)
	return result
}

func (c *Context) newValue(v interface{}) *JsValue {
	if c == nil || c.Raw == nil {
		return nil
	}
	val, err := v8go.NewValue(c.Raw.Isolate(), v)
	runtime.KeepAlive(v)
	if err != nil {
		panic(err)
	}
	return c.Wrap(val)
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
	args := make([]*Consumed, len(values))
	for i, v := range values {
		args[i] = c.NewString(v).Consume()
	}
	return c.NewArray(args...)
}
func (c *Context) NewArray(values ...*Consumed) *JsValue {
	a := c.RunScript("Array", "array")
	result := a.Call(a, values...)
	a.Release()
	return result
}
func (c *Context) NewObject() *JsValue {
	obj, err := c.objectTemplate.NewInstance(c.Raw)
	if err != nil {
		panic(err)
	}
	result := c.Wrap(obj.Value) //?
	return result
}
func (c *Context) NewArrayBuffer(data []byte) *JsValue {
	v := c.RunScript(fmt.Sprintf("new ArrayBuffer(%d)", len(data)), "arraybuffer.js")
	v8go.WriteToArrayBuffer(v.export(), data)
	return v
}
func (c *Context) NewFunctionTemplate(callback FunctionCallback) *FunctionTemplate {
	return newFunctionTemplate(c, callback)
}
func (c *Context) RunScript(script string, name string) *JsValue {
	result, err := c.Raw.RunScript(script, name)
	if err != nil {
		panic(err)
	}
	return c.Wrap(result)
}
func (c *Context) NullValue() *JsValue {
	return c.nullvalue
}

type Reusable struct {
	value *JsValue
}

func (r *Reusable) Consume() *Consumed {
	return &Consumed{
		JsValue:   r.value,
		noRelease: true,
	}
}
func (r *Reusable) FinalConsume() *Consumed {
	return &Consumed{
		JsValue: r.value,
	}
}
func (r *Reusable) Dispose() {
	r.value.Release()
}

// Consumed Should be used to wrap a JsValue that is consumed by a function call.
// You should not use JsValue directly after calling a function that consumes it.
type Consumed struct {
	*JsValue
	noRelease bool // if true, JsValue.Release() will not be called automatically
}

func (c *Consumed) Release() {
	if c.noRelease {
		return
	}
	c.JsValue.Release()
}

// JsValue is a wrapper for v8go.Value
//
// YOU MUST CALL Release() ON CONSUMED METHOD WHEN YOU FINISH USING IT.
//
// OR YOU CAN CALL CONSUM() METHOD TO GET A CONSUMED VALUE AND PASS IT TO FUNCTION CALL.
type JsValue struct {
	raw *v8go.Value
	ctx *Context
}

func (v *JsValue) Consume() *Consumed {
	return &Consumed{
		JsValue: v,
	}
}

// Reuse returns a reusable value that can be used to avoid unnecessary allocations.
// You should call Dispoe() or FinalConsume() on the returned value when you finish using it.
// You should not call Consume() on JsValue if you call ConsumeReuseble() on it.
func (v *JsValue) ConsumeReuseble() *Reusable {
	return &Reusable{
		value: v,
	}
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
func (v *JsValue) Call(recvr *JsValue, args ...*Consumed) *JsValue {
	if v.raw == nil {
		return nil
	}
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
	for i := range args {
		args[i].Release()
	}
	result := v.ctx.Wrap(val)
	return result
}

func (v *JsValue) Release() {
	if !v.raw.IsNullOrUndefined() {
		v.raw.Release()
	}
}

func (v *JsValue) Array() []*JsValue {
	result := []*JsValue{}
	length := v.Get("length")
	defer length.Release()
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
		item.Release()
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
func (v *JsValue) ArrayBufferContent() []byte {
	result := v8go.ArrayBufferContent(v.export())
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
func (v *JsValue) IsArrayBuffer() bool {
	result := v.export().IsArrayBuffer()
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

func (v *JsValue) MethodCall(methodName string, args ...*Consumed) *JsValue {
	fn := v.Get(methodName) // ensure method exists
	defer fn.Release()
	result := fn.Call(v, args...)
	runtime.KeepAlive(fn)
	runtime.KeepAlive(args)
	return result
}
func (v *JsValue) SetObjectMethod(ctx *Context, name string, fn FunctionCallback) {
	f := ctx.NewFunctionTemplate(fn).GetFunction(ctx).Consume()
	v.Set(name, f)
	runtime.KeepAlive(f)
}
func (v *JsValue) Get(key string) *JsValue {
	val, err := mustAsObject(v.export()).Get(key)
	if err != nil {
		panic(err)
	}
	result := v.ctx.Wrap(val)
	runtime.KeepAlive(v)
	return result
}
func (v *JsValue) GetIdx(idx uint32) *JsValue {
	val, err := mustAsObject(v.export()).GetIdx(idx)
	if err != nil {
		panic(err)
	}
	result := v.ctx.Wrap(val)
	runtime.KeepAlive(v)
	return result

}

func (v *JsValue) Set(key string, val *Consumed) {
	defer val.Release()
	err := mustAsObject(v.export()).Set(key, val.export())
	if err != nil {
		panic(err)
	}
}

func (v *JsValue) SetIdx(idx uint32, val *Consumed) {
	defer val.Release()
	err := mustAsObject(v.export()).SetIdx(idx, val.export())
	if err != nil {
		panic(err)
	}
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
	return result
}

func (v *JsValue) DeleteIdx(idx uint32) bool {
	result := mustAsObject(v.export()).DeleteIdx(idx)
	return result
}

type callback struct {
	cb  FunctionCallback
	ctx *Context
}

func (c *callback) call(info *v8go.FunctionCallbackInfo) *v8go.Value {
	rawargs := info.Args()
	args := make([]*Consumed, len(rawargs))
	for k, v := range rawargs {
		args[k] = c.ctx.Wrap(v).ConsumeReuseble().Consume()
	}
	this := c.ctx.Wrap(info.This().Value).ConsumeReuseble().Consume()
	defer this.JsValue.Release()
	fi := NewFunctionCallbackInfo(c.ctx, this, args...)
	result := c.cb(fi)

	for k := range args {
		args[k].JsValue.Release()
	}
	if result == nil {
		return nil
	}
	output := result.export()
	return output

}

type FunctionCallback func(info *FunctionCallbackInfo) *Consumed

func (f FunctionCallback) newCallback(c *Context) *callback {
	return &callback{cb: f, ctx: c}
}

func NewFunctionCallbackInfo(ctx *Context, this *Consumed, args ...*Consumed) *FunctionCallbackInfo {
	return &FunctionCallbackInfo{
		ctx:  ctx,
		this: this,
		args: args,
	}
}

type FunctionCallbackInfo struct {
	ctx  *Context
	args []*Consumed
	this *Consumed
}

func (i *FunctionCallbackInfo) Context() *Context {
	return i.ctx
}

// This returns the receiver object "this".
func (i *FunctionCallbackInfo) This() *Consumed {
	return i.this
}

// Args returns a slice of the value arguments that are passed to the JS function.
func (i *FunctionCallbackInfo) Args() []*Consumed {
	return i.args
}
func (i *FunctionCallbackInfo) GetArg(idx int) *Consumed {
	if idx < 0 || idx >= len(i.args) {
		return i.ctx.NullValue().Consume()
	}
	return i.args[idx]
}

type FunctionTemplate struct {
	tmpl *v8go.FunctionTemplate
}

func (t *FunctionTemplate) GetFunction(ctx *Context) *JsValue {
	fn := t.tmpl.GetFunction(ctx.Raw)
	return ctx.Wrap(fn.Value)
}
func newFunctionTemplate(c *Context, callback FunctionCallback) *FunctionTemplate {
	tmpl := v8go.NewFunctionTemplate(c.Raw.Isolate(), callback.newCallback(c).call)
	return &FunctionTemplate{
		tmpl: tmpl,
	}
}
func ExportRawValue(v *JsValue, noRelease bool) *v8go.Value {
	return v.raw
}
