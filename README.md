# gomux1

A very simple REST API server in GoLang using [gorilla/mux](https://github.com/gorilla/mux)

## Configuration
The app has a default config type `Config.go` in the `config` package. This type has certain defaults assigned to its fields, as well as environment variables tagged as `env:` that can be passed in to override those defaults. Please take a look at [config/config.go](config/config.go) for the fields, environment variable names and defaults.


### Environment variables
Any environment variable listed in [config/config.go](config/config.go) can be passed to the app executable via the standard Linux mechanism:

Example:
To override the server port:
```
SERVER_PORT=9090 ./gomux1
```

## Build
Standard mechanisms for GoLang build.
```
go build .
```

If the `go.sum` file is missing or you've updated `go.mod`:
```
go get github.com/rakhbari/gomux1
```

## Run
The app starts up 2 HTTP servers:
* A standard HTTP server
* A TLS-enabled HTTPS server

As such, the app requires the TLS certificate and key in one of 2 ways:

* `SERVER_TLS_CERT_PATH` set to a TLS cert bundle file - This file must be a bundle of the "leaf" cert + any CA intermediary certs + the CA root cert, all concatenated into a single file. Example:

```
SERVER_TLS_CERT_PATH="../openssl-cert/ca_chain-bundle.crt" SERVER_TLS_KEY_PATH="../openssl-cert/ca_intermediate_unencrypted.key" ./gomux1
```

* `SERVER_TLS_CERT_PATH` set to a TLS "leaf" cert file and `SERVER_TLS_CA_PATHS` set to a comma-delimited list of any CA intermediary certs and the CA root cert.

```
SERVER_TLS_CERT_PATH="../openssl-cert/leaf.crt" SERVER_TLS_KEY_PATH="../openssl-cert/ca_intermediate_unencrypted.key" SERVER_TLS_CA_PATHS="../openssl-cert/ca_intermediate.crt,../openssl-cert/ca_root.crt" ./gomux1
```

## Docker build/run
There is a Docker file in the repo which will build & run the app.

* Docker build:
```
docker build -t akcn/gomux1:latest .
```

* Docker run with only the HTTP server on the default port `8080`:
```
docker run --rm -p 8080:8080 akcn/gomux1:latest
```

* Docker run with a TLS bundle chain cert file:
```
docker run --rm -p 8080:8080 -p 8443:8443 \
  -v ${PWD}/../openssl-cert/ca_chain-bundle.crt:/cert.pem \
  -v ${PWD}/../openssl-cert/ca_intermediate_unencrypted.key:/cert.key \
  -e SERVER_TLS_CERT_PATH="/cert.pem" \
  -e SERVER_TLS_KEY_PATH="/cert.key" \
  akcn/gomux1:latest
```

* Docker run with a TLS "leaf" cert file + a intermediary CA certs and CA root cert files:
```
docker run --rm -p 8080:8080 -p 8443:8443 \
  -v ${PWD}/../openssl-cert/leaf.crt:/cert.pem \
  -v ${PWD}/../openssl-cert/ca_intermediate.crt:/ca_cert1.pem \
  -v ${PWD}/../openssl-cert/ca_root.crt:/ca_cert2.pem \
  -v ${PWD}/../openssl-cert/ca_intermediate_unencrypted.key:/cert.key \
  -e SERVER_TLS_CERT_PATH="/cert.pem" \
  -e SERVER_TLS_KEY_PATH="/cert.key" \
  -e SERVER_TLS_CA_PATHS="/ca_cert1.pem,/ca_cert2.pem" \
  akcn/gomux1:latest
```

## Test
As this is a very basic example app, the tests in `gomux1_test.go` don't do any extensive testing other than record the `content-type` and `status` code of the endpoints. But to run the tests in verbose mode:
```
go test -v
```

## Endpoints
2 endpoints are currently coded:
1. `ping`: Responds with a payload object of `response: pong!`
1. `health`: Responds with a payload object of `healthy: true`

## Standard Responses
Responses to all endpoints will be of the this standard structure, with the only difference being in what's contained in the `payload` field, which will vary depending on the endpoint hit.
```
{
  "requestId": "4a637cb1-f067-463d-94fe-ef51d392174c",
  "timestamp": "2022-03-22 13:27:00.4994833 -0700 PDT m=+8.117879601",
  "payload": <various-payloads-based-on-endpoint>
}
```
