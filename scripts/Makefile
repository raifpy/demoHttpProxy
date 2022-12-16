build:
	env CGO_ENABLED=0 env go build -ldflags="-w -s" --trimpath -v -x -o ./demo ../cmd/demo
	docker-compose up -d inf
	sleep 2 
	docker-compose build
	docker-compose down
	echo "docker-compose up to run"

run:
	docker-compose up