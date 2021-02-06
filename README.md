# go-escher-proxy

[![CI](https://github.com/zerosuxx/go-escher-proxy/workflows/CI/badge.svg)](https://github.com/zerosuxx/go-escher-proxy/actions?query=workflow%3ACI)

## Install
```
make build
```

## Run
```
make run
```

## Usage
```
curl -x localhost:8181 http://api.emarsys.net # http forced to https by default
```

## Build
```
make build
```

## Show available arguments
```
proxy -h
```

## Config (.proxy-config.json)
```
{"keyDB": [{"host": "api.emarsys.net", "accessKeyId": "app_suite_v1","apiSecret": "secret", "credentialScope": "eu/suite/ems_request"}]}
```
