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
