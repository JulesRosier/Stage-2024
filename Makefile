docker-build:
	docker build . -t  ghcr.io/julesrosier/stage-2024:latest --build-arg GIT_COMMIT=$$(git log -1 --format=%h)

docker-push:
	docker push ghcr.io/julesrosier/stage-2024:latest

docker-update:
	@make --no-print-directory docker-build
	@make --no-print-directory docker-push

build:
	go build -o ./tmp/main.exe ./cmd/master_event_gen

start:
	@./tmp/main.exe