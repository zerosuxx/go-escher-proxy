# go-escher-proxy

[![CI](https://github.com/zerosuxx/go-escher-proxy/workflows/CI/badge.svg)](https://github.com/zerosuxx/go-escher-proxy/actions?query=workflow%3ACI)

## Install
```
make install
```

## Run
```
make run
```

## Usage
```
curl -x localhost:8181 http://api.emarsys.net # http forced to https by default
curl -x localhost:8181 -H "X-Disable-Force-Https: 1" http://api.emarsys.net # http not forced to https
curl -H 'X-Target-Url: https://api.emarsys.net' http://localhost:8181
```

## Build
```
make build
```

## Show available arguments
```
proxy -h
```

## Config (proxy-config.json)
```
{
  "sites": {
    "api.emarsys.net": {
      "escherCredentials": {
        "disableBodyCheck": true,
        "accessKeyId": "app_suite_v1",
        "apiSecret": "dummySecret",
        "credentialScope": "eu/suite/ems_request"
      }
    }
  }
}
```
