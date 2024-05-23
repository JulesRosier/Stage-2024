# Stage-2024

This is the event generation part of Jules Rosier and Michiel Vereecke's internship assignment 2024.

## Dev

### Code gen

```sh
go install github.com/bufbuild/buf/cmd/buf@v1.30.0
buf generate
```

### Docker

#### Building

```sh
docker build . -t  ghcr.io/julesrosier/stage-2024:latest --build-arg GIT_COMMIT=$(git log -1 --format=%h)
docker push ghcr.io/julesrosier/stage-2024:latest
```

## config/config.yaml

```yaml
# Example config
logger:
  level: INFO
database:
  user: postgres
  password: password
  database: gen_oltp
  host: localhost
  port: 5432
kafka:
  brokers:
    - localhost:19092
  schemaRegistry:
    urls:
      - localhost:18081
  # Option for grouping of topics
  topicGrouping: true
```

## Dev Setup

### Docker

#### Setup

```sh
docker compose -f .\docker-compose-dev.yaml up -d
```

Make sure your `.env` is configured

```sh
docker compose up -d
```

### Dependencies

Latest version of Go

#### Building

```sh
make docker-build
```

### Running local

see `Makefile`

```sh
make start
```
