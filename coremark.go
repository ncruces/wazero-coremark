package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

//go:embed coremark-minimal.wasm
var binary []byte

var (
	module  api.Module
	runtime wazero.Runtime
)

func main() {
	boot := time.Now()
	ctx := context.Background()

	fmt.Println("Loading WebAssembly...")

	runtime = wazero.NewRuntime(ctx)
	env := runtime.NewHostModuleBuilder("env").NewFunctionBuilder().
		WithGoFunction(api.GoFunc(func(ctx context.Context, stack []uint64) {
			stack[0] = uint64(time.Since(boot) / time.Millisecond)
		}), nil, []api.ValueType{api.ValueTypeI64}).
		Export("clock_ms")
	_, err := env.Instantiate(ctx)
	if err != nil {
		log.Fatal(err)
	}

	module, err = runtime.Instantiate(ctx, binary)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Running CoreMark 1.0... [should take 12..20 seconds]")

	out, err := module.ExportedFunction("run").Call(ctx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Result: %.3f\n", math.Float32frombits(uint32(out[0])))
}
