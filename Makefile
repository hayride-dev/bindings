.PHONY: gen-ml

default: gen-ml

gen-ml:
	wit-bindgen-go generate --world imports --out ./go/ml/gen/ ./coven/ai/wit

gen-http: 
	wit-bindgen-go generate --world client --out ./go/http/gen/ ./coven/http/wit

gen-morph: 
	wit-bindgen-go generate --world imports --out ./go/morph/gen/ ./coven/morph/wit