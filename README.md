# slogx-go
  
  [![Go Report Card](https://goreportcard.com/badge/github.com/Evernorth/slogx)](https://goreportcard.com/report/github.com/Evernorth/slogx)
  [![GoDoc](https://godoc.org/github.com/Evernorth/slogx?status.svg)](https://godoc.org/github.com/Evernorth/slogx)
  [![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
  [![Release](https://img.shields.io/github/v/release/Evernorth/slogx)](https://github.com/Evernorth/slogx-go/releases)

## Description

A collection of `slog` extensions.
* `ContextHandler` allows you to add `slog` attributes (`slog.Attr` instances) to a `context.Context`.  These attributes are added to log records when the `*Context` function variants (`InfoContext`, `ErrorContext`, etc) on the logger are used.
* `LoggerBuilder` provides a simple way to build `slog.Logger` instances.
* `LevelManager` provides a way to manage `slog.LevelVar` instances from environment variables.

## Installation

``` go get -u github.com/Evernorth/slogx-go ```

## Usage

## Dependencies
See the [go.mod](go.mod) file.

## Support
If you have questions, concerns, bug reports, etc. See [CONTRIBUTING](CONTRIBUTING.md).

## License
slogx is open source software released under the [Apache 2.0 license](https://www.apache.org/licenses/LICENSE-2.0.html).

## Original Contributors
- Steve Sefton, Evernorth
- Shellee Stewart, Evernorth
- Neil Powell, Evernorth