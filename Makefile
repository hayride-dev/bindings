.PHONY: gen-ml gen-http gen-morph gen-silo gen-core gen

default: gen

gen-http: 
	wit-bindgen-go generate --world client-server --out ./go/wasihttp/gen/ coven/http/wit

gen-silo: 
	wit-bindgen-go generate --world imports --out ./go/silo/gen/ coven/silo/wit

gen-socket: 
	wit-bindgen-go generate --world exports --out ./go/socket/gen/ coven/socket/wit

gen-ai-exports:
	wit-bindgen-go generate --world exports --out ./go/ai/gen/exports coven/ai/wit

gen-ai-imports:
	wit-bindgen-go generate --world imports --out ./go/ai/gen/imports coven/ai/wit

gen-wasip2-io: 
	wit-bindgen-go generate --world io --out ./go/io/gen/io coven/wasip2/wit

gen: gen-http gen-silo gen-socket gen-ai-exports gen-ai-imports gen-wasip2-io

