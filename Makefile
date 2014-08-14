.PHONY: all clean

PROGS = metro server

all .DEFAULT: $(PROGS)

metro:
	go build metro.go

server:
	go build server.go

clean:
	rm -rf $(PROGS) *~ *.o *.a
