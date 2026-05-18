module github.com/vertex-language/vertex/cmd/vertex

go 1.22

require (
	github.com/vertex-language/vertex v0.0.0
	github.com/vertex-language/wasm-compiler v0.0.0-20260516130134-e90023166e7b
)

require (
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
)

replace github.com/vertex-language/vertex => ../../

replace github.com/vertex-language/wasm-compiler => /home/marss6414/wasm-compiler
