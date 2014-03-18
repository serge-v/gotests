GOPATH=$(HOME)/src/gocode
GO=/usr/local/go/bin/go

TARGETS=parse-dwml http-server http-client ftp-client

all: $(TARGETS)

%: %.go
	GOPATH=$(GOPATH) $(GO) build $<

getdeps:
	GOPATH=$(GOPATH) $(GO) get github.com/jlaffaye/ftp
	
%.run: %.go
	$(GO) run $<

clean:
	rm $(TARGETS)

help:
	$(GO) help
