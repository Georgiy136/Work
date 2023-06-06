
run: 
	go run cmd/app/main.go

docker_run:
	docker-compose up -d --build 

migrate:
	migrate -path migrations/ -database "postgres://postgres:4650@localhost/database?sslmode=disable" -verbose up

migrate-down:
	migrate -path migrations/ -database "postgres://postgres:4650@localhost/database?sslmode=disable" -verbose down

swag: swag_fmt
	swag init -d ./cmd/app,./ 

swag_fmt:
	swag fmt