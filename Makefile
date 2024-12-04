.PHONY: gen-ml gen-http gen-morph gen-silo gen-core gen

default: gen

gen-http: 
	wit-bindgen-go generate --world client-server --out ./go/wasihttp/gen/ ./coven/http/wit

gen-silo: 
	wit-bindgen-go generate --world imports --out ./go/silo/gen/ ./coven/silo/wit

gen: gen-http gen-silo 

