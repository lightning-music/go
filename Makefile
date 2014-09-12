.PHONY: install build test
PACKAGES := seq api binding
TESTS := seq

install .DEFAULT:
	$(foreach p, $(PACKAGES), (cd $(p) && go install);)

build:
	$(foreach p, $(PACKAGES), (cd $(p) && go build);)

# TODO: add more tests!
test:
	$(foreach t, $(TESTS), (cd $(t) && go test);)
