.PHONY: build

build:
	go build -o build/iso-edit cmd/main.go

clean:
	rm -rf build/*
	rm -rf /tmp/iso-test

run: clean build
	./build/iso-edit
