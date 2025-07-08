.PHONY: gen-ml gen-http gen-morph gen-silo gen-core gen

default: gen

# Note: All generated code lives under ./go/internal/gen/ this could cause conflicts.
gen-hayride:
	wit-bindgen-go generate --world hayride:bindings/hayride --out ./go/internal/gen/ ./wit
gen-hayride-x:
	wit-bindgen-go generate --world hayride:bindings/hayride-x --out ./go/internal/gen/ ./wit
gen-wasi:
	wit-bindgen-go generate --world hayride:bindings/wasi --out ./go/internal/gen/ ./wit

gen: gen-hayride gen-hayride-x gen-wasi

	