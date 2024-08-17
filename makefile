docker-compose-db:
	docker compose exec -it db psql -U root -d go-mailing

migrate-up:
	go run pkg/migration/cli/cli.go -path=internal/app/migrations -action=migrate-up

migrate-down:
	go run pkg/migration/cli/cli.go -path=internal/app/migrations -action=migrate-down

migrate-up-by-number:
	go run pkg/migration/cli/cli.go -path=internal/app/migrations -action=migrate-up-by-number -number=$(NUMBER)

migrate-down-by-number:
	go run pkg/migration/cli/cli.go -path=internal/app/migrations -action=migrate-down-by-number -number=$(NUMBER)

current-version:
	go run pkg/migration/cli/cli.go -path=internal/app/migrations -action=current-version