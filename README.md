<h1 align="center">dexus</h1>
<p align="center">
	<img src="https://github.com/mehmetumit/dexus/actions/workflows/build-test.yaml/badge.svg"/>
	<a href="https://goreportcard.com/report/github.com/mehmetumit/dexus"><img src="https://goreportcard.com/badge/github.com/mehmetumit/dexus"/></a>
	<a href="https://codecov.io/gh/mehmetumit/dexus"><img src="https://img.shields.io/codecov/c/github/mehmetumit/dexus/master.svg"/></a>
	<img src="https://github.com/mehmetumit/dexus/actions/workflows/release.yaml/badge.svg"/>
	<a href="https://github.com/mehmetumit/dexus/stargazers"><img src="https://img.shields.io/github/stars/mehmetumit/dexus?color=yellow" alt="Stars"/></a>
	<a href="https://github.com/mehmetumit/dexus/blob/main/LICENSE"><img src="https://img.shields.io/github/license/mehmetumit/dexus" alt="License"/></a>
</p>

---

## About

Dexus is a URL shortener that provides you with the ability to extend it dynamically according to your needs. You can utilize existing [adapters](./internal/adapters) or implement your own adapters based on [ports](./internal/core/ports). Afterwards, you only need to inject these adapters in [cmd/main.go](./cmd/main.go). Contributions of any new adapters are welcome.

## Usage
### Using binary release
Precompiled binaries are available on the [relases](https://github.com/mehmetumit/dexus/releases) page.
### Local Usage
```sh
# Clone the repository
git clone https://github.com/mehmetumit/dexus.git
# Change directory
cd dexus
# Build and run with make
make exec
# Or just build as "./build/dexus"
make build
```
### Using Docker
```sh
# Clone the repository
git clone https://github.com/mehmetumit/dexus.git
# Change directory
cd dexus
# Build image with make
make docker-build
# Run image with make
make docker-run
# Or manually
docker run -p 8080:8080 dexus:latest
```
### Development
```sh
# Live reload using 'entr'
make live-reload
```
```sh
# Run tests and get coverage
make test
```
```sh
# Run tests and get coverage with html output
make test-coverage-html
```
## Configuration

| Env Vars             | Default Values             |
|----------------------|----------------------------|
| `DEBUG_LEVEL`        | `true`                     |
| `HOST`               | `all interfaces`           |
| `PORT`               | `8080`                     |
| `READ_TIMEOUT_MS`    | `1000 ms`                  |
| `WRITE_TIMEOUT_MS`   | `1000 ms`                  |
| `CACHE_TTL_SEC`      | `60 s`                     |
| `YAML_REDIRECT_PATH` | `configs/redirection.yaml` |

## Contributing
Greatly appreciate contributions from the community. Your input, whether it's through code, documentation, bug reports, or suggestions, helps improve this project and make it even better.
## License
This project is licensed under the [GPLv3](./LICENSE) License.
