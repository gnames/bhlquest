# BHLquest

`BHLquest` is an AI app to query content of [Biodiversity Heritage Library]

## Installation

Make sure Go version on the computer is recent enough.

1. Install and run llmutil RESTful interface

Follow instructions at [llmutil documentation]

2. Install Go tools for `bhlquest`

Go to the root of `bhlquest` project and run:

```bash
make tools
go mod tidy
```

3. Install bhlquest

```bash
make install
```

## Usage

1. Install prerequisits and create database.

* Create `bhlnames` and `bhlquest` databases in PostgreSQL
* Download [bhlnames database dump]
* Restore `bhlnames`
* Install BHL text corpus

```bash
zstd -dc bhlnames-xxxx-xx-xx.zst|pg_restore -d bhlnames
```
* Install [pgvector extension]. It is needed to store vector data.

2. First run

First run wil create configuration file. Usually it is written at
`$HOME/.config/bhlquest.yaml`:

```bash
bhlquest -V
```

3. Edit `bhlquest.yaml`

Make sure that database settings, llmutil URL, directory which points to
BHL corpus settings reflect their actual values.

4. Import a subset of BHL data, and embed it into vectors.

```bash
bhlquest init --taxa 'Aves,Mammalia'
```

It will take a while!

5. Use `bhlquest` via command line or via RESTful API

### Command line Usage

```bash
bhlquest ask "What are the ecological niches of the Indigo Bunting?"
```

### Start RESTful API

```bash
bhlquest serve
```

With default settings [API description] should be accessible.

## Development

### Autogenerate API documentation

Install swag from the root of the project:

```bash
make tools
```

Run swag with the following command:

```bash
swag init -g server.go -d ./internal/io/web
```

Run the same command every time docs are updated.

[llmutil documentation]: https://github.com/gnames/llmutil
[pgvector extension]: https://github.com/pgvector/pgvector
[bhlnames database dump]: http://opendata.globalnames.org/dumps/bhlnames-2023-11-15.zst 
[API description]: http://0.0.0.0:8555/apidoc/
