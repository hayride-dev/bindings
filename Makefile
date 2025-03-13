.PHONY: gen-ml gen-http gen-morph gen-silo gen-core gen

default: gen

gen-imports:
	wit-bindgen-go generate --world hayride:bindings/imports --out ./go/internal/gen/imports ./wit

gen-exports:
	wit-bindgen-go generate --world hayride:bindings/exports --out ./go/internal/gen/exports ./wit

gen: gen-imports gen-exports

