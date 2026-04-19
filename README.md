# @jsonic/jsonc

This plugin allows the [Jsonic](https://jsonic.senecajs.org) JSON parser
to parse [JSONC](https://github.com/microsoft/node-jsonc-parser) format
files (JSON with Comments).

JSONC is a strict superset of JSON that adds single-line (`//`) and
block (`/* */`) comments. Trailing commas in objects and arrays can be
optionally enabled.

[![npm version](https://img.shields.io/npm/v/@jsonic/jsonc.svg)](https://npmjs.com/package/@jsonic/jsonc)
[![build](https://github.com/jsonicjs/jsonc/actions/workflows/build.yml/badge.svg)](https://github.com/jsonicjs/jsonc/actions/workflows/build.yml)
[![Coverage Status](https://coveralls.io/repos/github/jsonicjs/jsonc/badge.svg?branch=main)](https://coveralls.io/github/jsonicjs/jsonc?branch=main)
[![Known Vulnerabilities](https://snyk.io/test/github/jsonicjs/jsonc/badge.svg)](https://snyk.io/test/github/jsonicjs/jsonc)


| ![Voxgig](https://www.voxgig.com/res/img/vgt01r.png) | This open source module is sponsored and supported by [Voxgig](https://www.voxgig.com). |
| ---------------------------------------------------- | --------------------------------------------------------------------------------------- |


## Features

- Single-line comments: `// comment`
- Block comments: `/* comment */`
- Optional trailing commas in objects and arrays
- Strict JSON value parsing (no unquoted strings or hex numbers)
- Available in both TypeScript/JavaScript and Go


## TypeScript

### Install

```bash
npm install @jsonic/jsonc @jsonic/jsonic-next
```

### Quick Start

```typescript
import { Jsonic } from '@jsonic/jsonic-next'
import { Jsonc } from '@jsonic/jsonc'

const j = Jsonic.make().use(Jsonc)

// Parse JSONC with comments
const result = j('{ "name": "app", /* version */ "version": "1.0" }')
// => { name: "app", version: "1.0" }

// Enable trailing commas
const jc = Jsonic.make().use(Jsonc, { allowTrailingComma: true })
const config = jc('{ "debug": true, "verbose": false, }')
// => { debug: true, verbose: false }
```

### Options

| Option | Type | Default | Description |
|--------|------|---------|-------------|
| `allowTrailingComma` | `boolean` | `false` | Allow trailing commas in objects and arrays |
| `disallowComments` | `boolean` | `false` | Disable comment parsing (strict JSON mode) |


## Go

### Install

```bash
go get github.com/jsonicjs/jsonc/go
```

### Quick Start

```go
package main

import (
    "fmt"
    jsonc "github.com/jsonicjs/jsonc/go"
)

func main() {
    // Parse JSONC with comments
    result, err := jsonc.Parse(`{ "name": "app", /* version */ "version": "1.0" }`)
    if err != nil {
        panic(err)
    }
    fmt.Println(result)
    // => map[name:app version:1.0]

    // Enable trailing commas
    result, err = jsonc.Parse(
        `{ "debug": true, "verbose": false, }`,
        jsonc.JsoncOptions{AllowTrailingComma: boolPtr(true)},
    )
    fmt.Println(result)
    // => map[debug:true verbose:false]
}

func boolPtr(b bool) *bool { return &b }
```

### API

#### `Parse(src string, opts ...JsoncOptions) (any, error)`

Parse a JSONC string and return the result. Returns `map[string]any` for
objects, `[]any` for arrays, `float64` for numbers, `string`, `bool`,
or `nil`.

#### `MakeJsonic(opts ...JsoncOptions) *jsonic.Jsonic`

Create a configured jsonic instance for JSONC parsing. Use this when you
need to parse multiple inputs with the same configuration.

#### `JsoncOptions`

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `AllowTrailingComma` | `*bool` | `false` | Allow trailing commas in objects and arrays |
| `DisallowComments` | `*bool` | `false` | Disable comment parsing |


## JSONC Format

JSONC follows [RFC 8259](https://tools.ietf.org/html/rfc8259) (JSON) with
these extensions:

- **Line comments**: `//` to end of line
- **Block comments**: `/* */` (non-nesting)
- **Trailing commas**: optional, in objects and arrays

All other rules follow standard JSON:
- Strings must be double-quoted
- Only standard escape sequences: `\"` `\\` `\/` `\b` `\f` `\n` `\r` `\t` `\uXXXX`
- Numbers: integer, decimal, scientific notation (no hex, octal, or binary)
- Keywords: `true`, `false`, `null` (case-sensitive)
- Property names must be double-quoted strings


## Acknowledgments

Conformance testing uses third-party corpora under MIT License:

- [nst/JSONTestSuite](https://github.com/nst/JSONTestSuite) by Nicolas
  Seriot — vendored as a git submodule at `test/JSONTestSuite/`.
- [microsoft/node-jsonc-parser](https://github.com/microsoft/node-jsonc-parser) —
  parse-level test cases ported into `test/jsonc.test.ts`.

See [THIRD_PARTY_NOTICES.md](./THIRD_PARTY_NOTICES.md) for details.


## License

MIT. Copyright (c) 2021-2025 Richard Rodger and contributors.
