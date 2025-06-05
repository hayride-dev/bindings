.PHONY: gen-ml gen-http gen-morph gen-silo gen-core gen

default: gen

gen-imports:
	wit-bindgen-go generate --world hayride:bindings/sdk --out ./go/internal/gen ./wit

gen: gen-imports

