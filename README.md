# BHLquest

`BHLquest` is an AI application designed to query the content of
 the [Biodiversity Heritage Library]

## Installation

Ensure that your computer's Go version is up to date.

1. Install and Run `llmutil` RESTful Interface

Follow instructions at [llmutil documentation].

2. Install Go Tools for `bhlquest`

Navigate to the root of the `bhlquest` project and execute:

```bash
make tools
go mod tidy
```

3. Install `bhlquest`

```bash
make install
```

## Usage

1. Install Prerequisits and Create Database

* Create `bhlnames` and `bhlquest` databases in PostgreSQL.
* Download the [bhlnames database dump].
* Restore the `bhlnames` database:

```bash
zstd -dc bhlnames-xxxx-xx-xx.zst|pg_restore -d bhlnames
```

* Install the BHL text corpus.
* Install [pgvector extension], which is necessary for storing vector data.

2. Initial Run

First run wil create a configuration file, typically located at
`$HOME/.config/bhlquest.yaml`:

```bash
bhlquest -V
```

3. Edit `bhlquest.yaml`

Ensure that database settings, `llmutil` URL, directory which points to
the BHL corpus settings accurately reflect their real values.

4. Import and Embed a Subset of BHL data into Vectors

```bash
bhlquest init --classes 'Aves,Mammalia' --rebuild-db

# if init process was interrupted, it can be continued from the place it
# stopped. Make sure the `--classes` values are the same as before.
# Remove `--rebuild-db` flag to keep alreday saved data.

bhlquest init --classes 'Aves,Mammalia'
```

If no clases are specified, the process will include all BHL data.

It will take a **while**!

5. Use `bhlquest` via Command Line or RESTful API

### Command line Usage

```bash
bhlquest ask "What are the ecological niches of the Indigo Bunting?"
```

### Starting RESTful API

```bash
bhlquest serve
```

With the default settings, the [API description] should be accessible.

## Development

### Auto-generating API Documentation

Install [`swag`] from the root of the project:

```bash
make tools
```

To update the documentation, execute:

```bash
swag init -g server.go -d ./internal/io/web
```

Run this command each time the docs are updated.

### Generating a Client for a Specific Language

You can use the [`openapi-generator`] to create an API client in a language of
your choice. For instance, to generate a Ruby client, run the following
command from the project's root:

```bash
openapi-generator generate -i ./docs/swagger.yaml -g ruby -o ~/tmp/bhlquest --additional-properties gemName=bhlquest
```


[Biodiversity Heritage Library]: https://www.biodiversitylibrary.org/
[llmutil documentation]: https://github.com/gnames/llmutil
[pgvector extension]: https://github.com/pgvector/pgvector
[bhlnames database dump]: http://opendata.globalnames.org/dumps/bhlnames-2023-11-15.zst 
[API description]: http://0.0.0.0:8555/apidoc/
[`swag`]: https://github.com/swaggo/swag
[`openapi-generator`]: https://github.com/OpenAPITools/openapi-generator
