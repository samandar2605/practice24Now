POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DATABASE=users

include .env
  
DB_URL=postgresql://postgres:12345@localhost:5432/users?sslmode=disable


print:
	echo "$(DB_URL)"
	
swag:
	swag init -g api/api.go -o api/docs

run:
	go run "/home/samandar/go/src/github.com/practice2311/cmd/main.go"


migrate_file:
	migrate create -ext sql -dir migrations/ -seq alter_some_table



migrateup:
	migrate -path migrations -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path migrations -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path migrations -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path migrations -database "$(DB_URL)" -verbose down 1

.PHONY: start migrateup migratedown