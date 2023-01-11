# Exporter

Exporter is a go project for exporting the psql table into clickhouse

## Installation

```bash
go mog tidy
```

## Usage

```go
go run cmd/main.go
```

## Notes
1. All env file examples can be found in config/config.go file
2. Please tune the docker-compose.yaml file based on the needs, because the followings
    - PSQL might have been isntalled locally - configure with external IP address if you want to reach to locally installed psql
    - Clickhouse might have been installed locally - configure with external IP address if you want to reach the locally installed psql
3. Go through Makefile for building images and pushing them into some repository

## Contributing

## License

[MIT](https://choosealicense.com/licenses/mit/)
