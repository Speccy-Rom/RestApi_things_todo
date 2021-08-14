build:
	docker-compose build restapi_things_todo

run:
	docker-compose up restapi_things_todo

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5436/postgres?sslmode=disable' up
