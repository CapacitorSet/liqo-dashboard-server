# liqo-dashboard-server

A REST API for Liqo.

## Interface

The endpoints are described in OpenAPI format in `api.yaml`.

## Hacking

To change the API methods, edit `api.yaml` and regenerate the code:

```sh
go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen
oapi-codegen -old-config-style -package api -generate types,server,spec api.yaml > api/api.go
oapi-codegen -old-config-style -package client -generate types,client api.yaml > client/client.go
```