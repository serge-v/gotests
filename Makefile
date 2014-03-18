GO=/usr/local/go/bin/go

TARGETS=parse-dwml http-server http-client

all: $(TARGETS)

%: %.go
	$(GO) build $<

%.run: %.go
	$(GO) run $<

clean:
	rm $(TARGETS)

help:
	$(GO) help
