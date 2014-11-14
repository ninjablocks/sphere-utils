all:
	scripts/build.sh

install:
	cd tools/sphere-go-serial && go install
	cd tools/sphere-go-config && go install

clean:
	rm -f bin/* || true
	rm -rf .gopath || true

test:
	go test -v ./...

vet:
	go vet ./...

.PHONY: all clean test vet
