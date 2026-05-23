EXECUTABLE = bin/server
PKG = ./cmd/server
DB_NAME = vms.db

build:
	go build -o $(EXECUTABLE) $(PKG)
	cp -r web/templates bin/

run:
	TEMPLATE_DIR="" $(EXECUTABLE)

dev:
	go run $(PKG)

init-db:
	sqlite3 $(DB_NAME) < ./db/init_db.sql

clean:
	rm -rf bin
	rm -f $(DB_NAME)

.PHONY: build run clean dev init-db
