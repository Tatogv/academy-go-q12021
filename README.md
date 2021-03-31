# Go Academy 2021

This project is a small API client developed in Go, with the purpose to test and demonstrate my skills in the language.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)

## Installation

Download to your project directory (make sure you have [Go](https://golang.org/doc/install)installed), then run:

```sh
go run .
```

## Usage



This API includes the following endpoints:

|Route|Description|
--- | ---
|/read|Read the contents of the stored CSV file.
|/read/{id}|Retrieve information of a specific line of the CSV.
|/getBerries|Retrieve a list of berries and save them on a CSV
|/readConcurrently|Read a CSV using Go's powerful concurrency.

## TO DO
- Create swagger file
- Implement Viper for configuration files

