.SILENT:
.PHONY: run migrate-tasks

run:
	docker-compose up --remove-orphans --build

migrate-tasks:
	migrate -path ./task-svc/db -database 'postgres://postgres:qwerty@0.0.0.0:5433/tasks?sslmode=disable' up