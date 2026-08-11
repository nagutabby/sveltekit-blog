.PHONY: generate
generate:
	cd proto && PATH="$$PWD/../web/node_modules/.bin:$$PATH" buf generate

.PHONY: db-up
db-up:
	docker compose up -d

.PHONY: db-down
db-down:
	docker compose down -v
