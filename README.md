# bindings

This repository contains generated bindings and language-specific wrappers for WebAssembly Interfaces defined in [Coven](https://github.com/hayride-dev/coven).

## WIT

Bindings are generated from the WIT files in the `wit` directory. 

In order to generate bindings, a world must be defined. It is not recommended to include our `bindings` worlds in your own projects. Instead, you should define you own worlds that include the necessary interfaces and use our bindings to satisfy WebAssembly imports/exports of a component.

### Updating WIT Dependencies

WIT Dependencies are managed using [wit-deps](https://github.com/bytecodealliance/wit-deps).

You can add dependencies to the `wit/deps.toml` and pull them in with:

```bash 
wit-deps update
```

## Go Bindings 

Golang bindings are generated using [wit-bindgen-go](https://github.com/bytecodealliance/go-modules). 

After installing required dependencies, you can generate the bindings via the Makefile:

```bash
make gen
```

## Future languages
We plan to support additional languages in the future. If you are interested in contributing bindings for a specific language, please open an issue or submit a pull request.

## Contributing
Contributions are welcome! If you'd like to contribute, please follow these steps:

- Fork the repository.
- Create a new branch for your feature or bug fix.
- Submit a pull request with a detailed description of your changes.

## License
This project is licensed under the MIT License. See the LICENSE file for details