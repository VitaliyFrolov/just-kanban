up-dev-app:
	docker-compose -f docker-compose.dev.yml up

down-dev-app:
	docker-compose -f docker-compose.dev.yml down

clear-dev-app:
	docker-compose -f docker-compose.dev.yml down -v