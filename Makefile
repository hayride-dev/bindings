.PHONY: gen-ml gen-http gen-morph gen-silo gen-core gen

default: gen

# Note: All generated code lives under ./go/internal/gen/ this could cause conflicts.
gen-imports:
	wit-bindgen-go generate --world hayride:bindings/imports --out ./go/internal/gen/imports ./wit
gen-exports:
	wit-bindgen-go generate --world hayride:bindings/exports --out ./go/internal/gen/exports ./wit
gen-wasi:
	wit-bindgen-go generate --world hayride:bindings/wasip2 --out ./go/internal/gen/wasip2 ./wit
gen-types:
	wit-bindgen-go generate --world hayride:bindings/types --out ./go/internal/gen/types ./wit

gen: gen-imports gen-exports gen-wasi gen-types