# gomux1

A very simple REST API server in GoLang using [gorilla/mux](https://github.com/gorilla/mux)

## Configuration
The app can have its configuration passed to it in one of 2 ways, or combination of both.

### `config.yaml` file
`config.yaml` will be read in to supply server `host` and `port` values.

### Environment variables
You can pass the following environment variables at app startup time, which will override anything that's been defined in `config.yaml` file.
* `SERVER_HOST`: Defaults to `localhost`
* `SERVER_PORT`: Defaults to `8080`

Example:
To override the server port:
```
SERVER_PORT=9090 ./gomux1
```

## Run
Very simple server startup
```
./gomux1
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
