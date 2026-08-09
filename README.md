# OpenStalk

![Test](https://github.com/ZoZo-182/OpenStalk/actions/workflows/test.yml/badge.svg)

OpenStalk is a Go CLI that helps developers discover active open source projects by
searching GitHub for recent open pull requests. It is built for finding projects people are actively working on, so you can discover
more tools or just cool projects in general.

## Features

- Search recent open pull requests on GitHub
- Filter by programming language
- Limit result count
- Output text or JSON
- Use `GITHUB_TOKEN` for higher GitHub API limits

## Installation
Install the latest tagged release with Go:

```sh
go install github.com/ZoZo-182/openstalk/v2@latest
```

Make sure your Go binary directory is on your `PATH`:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Usage

Search for recent open pull requests:

```sh
openstalk
```

Filter by language:

```sh
openstalk search --language go
```

Use short flags:

```sh
openstalk search -l go -d 7 -n 5
```

Show help:

```sh
openstalk --help
openstalk search --help
```

## GitHub API Token

OpenStalk can run without a GitHub token, but unauthenticated GitHub API requests are heavily rate-limited.

To increase the rate limit, create a GitHub personal access token and set:

```sh
export GITHUB_TOKEN=some_token
```

## License
MIT. See `LICENSE` for more information.
