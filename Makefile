.SILENT:
.PHONY: run migrate-tasks

run:
	docker-compose build --no-cache
	docker-compose up --remove-orphans

migrate-tasks:
	migrate -path ./task-svc/schema -database 'postgres://postgres:qwerty@0.0.0.0:5433/tasks?sslmode=disable' up

migrate-tasks-drop:
	migrate -path ./task-svc/schema -database 'postgres://postgres:qwerty@0.0.0.0:5433/tasks?sslmode=disable' drop