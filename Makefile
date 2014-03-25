GOPATH=$(HOME)/src/gocode
GO=/usr/local/go/bin/go

TARGETS=parse-dwml http-server http-client ftp-client mysql-client parse-json

all: $(TARGETS)

%: %.go
	GOPATH=$(GOPATH) $(GO) build $<

getdeps:
	GOPATH=$(GOPATH) $(GO) get github.com/jlaffaye/ftp
	GOPATH=$(GOPATH) $(GO) get github.com/go-sql-driver/mysql
	
%.run: %.go
	GOPATH=$(GOPATH) $(GO) run $<

clean:
	rm $(TARGETS)

help:
	$(GO) help
