.PHONY: install build test clean
PACKAGES := seq api binding
TESTS := seq
# GOINSTALL := go install -ldflags -w -gcflags "-N -l"
GOINSTALL := go install -a

install .DEFAULT:
	$(GOINSTALL) github.com/lightning/go/api
	$(GOINSTALL) github.com/lightning/go/api/handler
	$(GOINSTALL) github.com/lightning/go/binding
	$(GOINSTALL) github.com/lightning/go/seq
	$(GOINSTALL) github.com/lightning/go/types
	cd example/clients && go build osc-rand.go

build:
	$(foreach p, $(PACKAGES), (cd $(p) && go build);)

# TODO: add more tests!
test:
	$(foreach t, $(TESTS), (cd $(t) && go test);)

clean:
	rm -rf lightningd example/clients/osc-rand \
           examples/play-sample
