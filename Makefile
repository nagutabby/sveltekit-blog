.PHONY: generate
generate:
	cd proto && PATH="$$PWD/../web/node_modules/.bin:$$PATH" buf generate

.PHONY: db-migrate
db-migrate:
	cd backend && goose -dir db/migrations sqlite3 $${SQLITE_PATH:-backend.db} up
