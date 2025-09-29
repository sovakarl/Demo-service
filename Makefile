container_name=postgres_db

pathToApp=./cmd/app/main.go

default:docker

docker:down build up

logs:
	docker logs -f ${container_name}

build:
	docker-compose build 

up:
	docker-compose up -d 

stop:
	docker-compose stop 

down:
	docker-compose down -v 

shell:
	docker exec -it ${container_name} bash

app:
	go run /cmd/main.go