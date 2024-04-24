# Stage-2024

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
```
