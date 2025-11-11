DATABASE_URL=mysql://root:Halloexpress123@@tcp(localhost:3306)/erp_finpos?charset=utf8mb4&parseTime=True&loc=Local

dev:
	go run cmd/web/main.go

migrate-up:
	migrate -database '${DATABASE_URL}' -path db/migrations up

migrate-down:
	migrate -database '${DATABASE_URL}' -path db/migrations down

migrate-new:
	@if [ -z "$(n)" ]; then \
		echo "❌ Error: Harap berikan nama migrasi, contoh:"; \
		echo "   make migrate-new n=create_table_users"; \
		exit 1; \
	fi
	migrate create -ext sql -dir db/migrations $(n)
