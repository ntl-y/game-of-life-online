env GOOS=js GOARCH=wasm go build -o test.wasm test
cp $(go env GOROOT)/misc/wasm/wasm_exec.js .
