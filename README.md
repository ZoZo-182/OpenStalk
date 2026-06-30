![cute badge](https://github.com/ZoZo-182/OpenStalk/actions/workflows/test.yml/badge.svg)

## About The Project
OpenStalk is a Go CLI that helps developers discover active open source projects by
searching GitHub for recent open pull requests.

It is built for finding projects people are actively working on, so you can discover
more tools or just cool projects in general.

## Installation
Install the latest tagged release with Go:

```sh
go install github.com/ZoZo-182/openstalk@latest
```

Make sure your Go binary directory is on your `PATH`:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

Then run:

```sh
openstalk search
```

## Github API Token
OpenStalk can run without a Github token, but unauthenticated Github API requests are heavily rate-limited.

To increase the rate limit, create a Github personal access token and set:
```sh 
export GITHUB_TOKEN=some_token
```


## Cute Demo (OLD PLEASE UPDATE -> Usage)
![Demo](https://media1.giphy.com/media/v1.Y2lkPTc5MGI3NjExOWkydndtOXRvZHNjMTZydjl6amlwY2swMWh6Yno4eXF1cDZyZnNjOCZlcD12MV9pbnRlcm5hbF9naWZfYnlfaWQmY3Q9Zw/JPCjH5tkIps3xY3pFt/giphy.gif)


## License
MIT. See `LICENSE.txt` for more information.
