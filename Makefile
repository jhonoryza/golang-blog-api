build:
	docker build -t jhonoryza/api_blog .
run:
	docker compose up -d
update:
	git pull origin main && make build && make run