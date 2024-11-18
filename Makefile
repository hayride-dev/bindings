.PHONY: gen-ml gen-http gen-morph gen-silo gen-core gen

default: gen

gen-ml:
	wit-bindgen-go generate --world imports --out ./go/ml/gen/ ./coven/ai/wit

gen-http: 
	wit-bindgen-go generate --world client-server --out ./go/wasihttp/gen/ ./coven/http/wit

gen-morph: 
	wit-bindgen-go generate --world imports --out ./go/morph/gen/ ./coven/morph/wit

gen-silo: 
	wit-bindgen-go generate --world imports --out ./go/silo/gen/ ./coven/silo/wit

gen-core: 
	wit-bindgen-go generate --world imports --out ./go/core/gen/ ./coven/core/wit

gen:  gen-ml gen-http gen-morph gen-silo gen-core

