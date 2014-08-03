.PHONY: all clean

PROGS = metro

all .DEFAULT: $(PROGS)

metro:
	go build metro.go

clean:
	rm -rf $(PROGS) *~ *.o *.a
