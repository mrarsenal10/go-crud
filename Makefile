.PHONY: down_application
down_application:
	docker-compose -f deployment/docker-compose.yml down

.PHONY: build_application
build_application:
	make down_application

	docker-compose -f deployment/docker-compose.yml down
	docker-compose -f deployment/docker-compose.yml build
	docker-compose -f deployment/docker-compose.yml up -d