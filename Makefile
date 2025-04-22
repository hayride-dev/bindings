.PHONY: gen-ml gen-http gen-morph gen-silo gen-core gen

default: gen

gen-imports:
	wit-bindgen-go generate --world hayride:bindings/imports --out ./go/gen/imports ./wit

gen-exports:
	wit-bindgen-go generate --world hayride:bindings/exports --out ./go/gen/exports ./wit

gen: gen-imports gen-exports

