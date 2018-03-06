migrate:
	migrate -database=postgres://max@localhost:5432/dumbwall?sslmode=disable -path=etc/migrations/ up

migrate-force:
    migrate -database=postgres://max@localhost:5432/dumbwall?sslmode=disable -force 2
