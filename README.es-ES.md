# OpenStalk

![Test](https://github.com/ZoZo-182/OpenStalk/actions/workflows/test.yml/badge.svg)

OpenStalk es una CLI en Go que ayuda a los desarrolladores a descubrir proyectos de código abierto activos buscando en GitHub pull requests abiertos recientes. Está construida para encontrar proyectos en los que la gente está trabajando activamente, así puedes descubrir más herramientas o simplemente proyectos interesantes en general.

## Características

- Buscar pull requests abiertos recientes en GitHub
- Filtrar por lenguaje de programación
- Filtrar por rango de estrellas
- Limitar cantidad de resultados
- Salida en texto o JSON
- Usar `GITHUB_TOKEN` para mayores límites de API de GitHub

## Instalación
Instala la última versión publicada con Go:

```sh
go install github.com/ZoZo-182/openstalk@latest
```

Asegúrate de que tu directorio de binarios de Go esté en tu `PATH`:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Uso

Buscar pull requests abiertos recientes:

```sh
openstalk search
```

Filtrar por lenguaje:

```sh
openstalk search --language go
```

Usar banderas cortas:

```sh
openstalk search -l go -d 7 -n 5
```

Mostrar ayuda:

```sh
openstalk --help
openstalk search --help
```

## Token de API de GitHub

OpenStalk puede funcionar sin un token de GitHub, pero las solicitudes a la API de GitHub sin autenticación están muy limitadas.

Para aumentar el límite, crea un token de acceso personal de GitHub y configura:

```sh
export GITHUB_TOKEN=some_token
```

## Licencia
MIT. Consulta `LICENSE` para más información.
