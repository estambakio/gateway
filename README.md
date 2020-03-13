[![CircleCI](https://circleci.com/gh/estambakio/gateway.svg?style=shield)](https://circleci.com/gh/estambakio/gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/estambakio/gateway)](https://goreportcard.com/report/github.com/estambakio/gateway)
[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/estambakio/gateway)](https://cloud.docker.com/repository/docker/estambakio/gateway)
[![MicroBadger Size (tag)](https://img.shields.io/microbadger/image-size/estambakio/gateway/master)](https://cloud.docker.com/repository/docker/estambakio/gateway)
[![Docker Cloud Automated build](https://img.shields.io/docker/cloud/automated/estambakio/gateway)](https://cloud.docker.com/repository/docker/estambakio/gateway)
[![license](https://img.shields.io/github/license/estambakio/gateway.svg?style=flat-square)](LICENSE)

# gateway

## State

Paused. Primary use case (handling 'maintenance' case in Kubernetes cluster when admins have access to all regular ingress routes while other users receive 'maintenance in progress' page) can be solved via built-in functionality in SSO tool, e.g. Keycloak.

## Usage

Create config with proxy rules:

```
rules:
  - from: "/task/"
    to: "https://jsonplaceholder.typicode.com/todos/"
```

Starts gateway with this config on your local machine on port 3000:

```
docker run -it -p 3000:3000 \
  -v /path/to/config.yaml:/config.yaml \
  estambakio/gateway:master \
    --config /config.yaml
```

Open terminal and make a request:

```
curl -v http://localhost:3000/task/1
```

This will respond with content from `https://jsonplaceholder.typicode.com/todos/1`

### Args
- `-c or --config` - path to config file, **required**

### Environment variables
- `PORT` - port to start server on, 3000 by default

## Development

`go test -v -cover ./...`
