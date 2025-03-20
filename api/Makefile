init-api: generate-migrations migrations-up run-api

run-api:
	go run cmd/app/main.go

run-tests:
	go test ./...

generate-migrations:
	go run scripts/generate_migrations.go

# Being used for generating mock structures for testing complex structures
generate-interfaces-mock:
	mockgen -destination ./mocks/$(FILE_NAME).mock.go -package=mocks $(PACKAGE_PATH) $(INTERFACE)

migrations-up:
	go run cmd/migrations/main.go -direction=up

migrations-down:
	go run cmd/migrations/main.go -direction=down

migrations-drop:
	go run cmd/migrations/main.go -direction=drop