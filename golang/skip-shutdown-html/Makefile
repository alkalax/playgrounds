EXECUTABLE = bin/server
PKG = ./cmd/server

build:
	go build -o $(EXECUTABLE) $(PKG)
	cp -r web/templates bin/

run:
	TEMPLATE_DIR="" $(EXECUTABLE)

dev:
	go run $(PKG)

clean:
	rm -rf bin

.PHONY: build run clean dev