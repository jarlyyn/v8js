package v8plugin

import (
	"context"
	"testing"

	"github.com/herb-go/herbplugin"
)

var moduleinitoutput string
var modulebootoutput string
var modulecloseoutput string

func newTestModule(output string) *herbplugin.Module {
	return herbplugin.CreateModule(
		output,
		func(ctx context.Context, p herbplugin.Plugin, next func(ctx context.Context, plugin herbplugin.Plugin)) {
			moduleinitoutput += output
			next(ctx, p)
		},
		func(ctx context.Context, p herbplugin.Plugin, next func(ctx context.Context, plugin herbplugin.Plugin)) {
			modulebootoutput += output
			next(ctx, p)
		},
		func(ctx context.Context, p herbplugin.Plugin, next func(ctx context.Context, plugin herbplugin.Plugin)) {
			modulecloseoutput += output
			next(ctx, p)
		},
	)
}
func TestPlugin(t *testing.T) {
	p := New()
	if p.PluginType() != PluginType {
		t.Fatal(p)
	}
	moduleinitoutput = ""
	modulebootoutput = ""
	modulecloseoutput = ""
	i := NewInitializer()
	i.Modules = []*herbplugin.Module{
		newTestModule("test1"),
		newTestModule("test2"),
		newTestModule("test3"),
	}
	i.MustApplyInitializer(p)
	herbplugin.Lanuch(p, herbplugin.NewOptions())
	if moduleinitoutput != "test1test2test3" || modulebootoutput != "test1test2test3" || modulecloseoutput != "" {
		t.Fatal(moduleinitoutput, modulebootoutput, modulecloseoutput)
	}
	p.MustClosePlugin()
	if moduleinitoutput != "test1test2test3" || modulebootoutput != "test1test2test3" || modulecloseoutput != "test3test2test1" {
		t.Fatal(moduleinitoutput, modulebootoutput, modulecloseoutput)
	}
}
