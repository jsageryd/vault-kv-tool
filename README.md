# vault-kv-tool

[![Build Status](https://travis-ci.com/jsageryd/vault-kv-tool.svg?branch=master)](https://travis-ci.com/jsageryd/vault-kv-tool)
[![Go Report Card](https://goreportcard.com/badge/github.com/jsageryd/vault-kv-tool)](https://goreportcard.com/report/github.com/jsageryd/vault-kv-tool)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](https://github.com/jsageryd/vault-kv-tool#license)

This is a tool for importing/exporting KV (version 1) trees to/from
[Vault](https://www.vaultproject.io/).

## Installation
```
go get -u github.com/jsageryd/vault-kv-tool
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

## License
Copyright (c) 2019 Johan Sageryd <j@1616.se>

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
