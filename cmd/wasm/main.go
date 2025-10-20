package main

import (
	"bgweb-api/internal/api"
	"bgweb-api/internal/gnubg"
	"bgweb-api/internal/openapi"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
)

//go:embed data
var data embed.FS

var initialized bool

//export wasm_init
func wasm_init() {
	if initialized {
		return
	}

	// root embedded fs to data/
	dataDir, err := fs.Sub(data, "data")
	if err != nil {
		panic(err)
	}

	if err := gnubg.Init(dataDir); err != nil {
		panic(err)
	}

	initialized = true
}

//export wasm_get_moves
func wasm_get_moves(jsonArgs string) string {
	if !initialized {
		return `{"error": "WASM not initialized. Call wasm_init() first."}`
	}

	var args openapi.MoveArgs

	if err := json.Unmarshal([]byte(jsonArgs), &args); err != nil {
		return fmt.Sprintf("{\"error\": \"%v\"}", err.Error())
	}

	moves, err := api.GetMoves(args)

	if err != nil {
		return fmt.Sprintf("{\"error\": \"%v\"}", err.Error())
	}

	bytes, err := json.Marshal(moves)

	if err != nil {
		return fmt.Sprintf("{\"error\": \"%v\"}", err.Error())
	}

	return string(bytes)
}

func main() {
	// For React Native WebAssembly, main() should not block
	// The exported functions will be called directly by react-native-webassembly
}
