.PHONY: install build test clean
PACKAGES := seq api binding
TESTS := seq

install .DEFAULT:
	go install github.com/lightning/go/api
	go install github.com/lightning/go/api/handler
	go install github.com/lightning/go/binding
	go install github.com/lightning/go/seq
	go install github.com/lightning/go/types
	go build lightningd.go

build:
	$(foreach p, $(PACKAGES), (cd $(p) && go build);)

# TODO: add more tests!
test:
	$(foreach t, $(TESTS), (cd $(t) && go test);)

clean:
	rm lightningd
