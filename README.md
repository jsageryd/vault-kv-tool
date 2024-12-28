# vault-kv-tool

[![Build Status](https://github.com/jsageryd/vault-kv-tool/workflows/ci/badge.svg)](https://github.com/jsageryd/vault-kv-tool/actions?query=workflow%3Aci)
[![Go Report Card](https://goreportcard.com/badge/github.com/jsageryd/vault-kv-tool)](https://goreportcard.com/report/github.com/jsageryd/vault-kv-tool)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/jsageryd/vault-kv-tool#license)

This is a tool for importing/exporting KV (version 1) trees to/from
[Vault](https://www.vaultproject.io/).

## Installation
```
go install github.com/jsageryd/vault-kv-tool@latest
```

## Usage example
```
vault server -dev
```

```
$ export VAULT_ADDR=http://localhost:8200/
$ export VAULT_TOKEN=<your token>
$ vault secrets enable -version=1 -path=sec kv
$ vault write sec/foo/bar baz=qux
$ vault write sec/bar/baz qux=foo
$ vault-kv-tool -root=sec | jq . | tee data.json
{
  "bar/": {
    "baz": {
      "qux": "foo"
    }
  },
  "foo/": {
    "bar": {
      "baz": "qux"
    }
  }
}
$ vault secrets enable -version=1 -path=sec2 kv
$ vault-kv-tool -root=sec2 -write < data.json
$ vault-kv-tool -root=sec2 | jq .
{
  "bar/": {
    "baz": {
      "qux": "foo"
    }
  },
  "foo/": {
    "bar": {
      "baz": "qux"
    }
  }
}
$ vault read -format=json sec2/foo/bar | jq .data
{
  "baz": "qux"
}
```
