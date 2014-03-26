GOPATH=$(HOME)/src/gocode
GO=/usr/local/go/bin/go

TARGETS=parse-dwml http-server http-client ftp-client mysql-client parse-json cache-test

all: $(TARGETS)

version/version.go: *.go
	./gen-version.sh

%: %.go
	GOPATH=$(GOPATH) $(GO) build $<

http-server: version/version.go http-server.go http-server-config.go
	GOPATH=$(GOPATH) $(GO) build http-server.go http-server-config.go

getdeps:
	GOPATH=$(GOPATH) $(GO) get github.com/jlaffaye/ftp
	GOPATH=$(GOPATH) $(GO) get github.com/go-sql-driver/mysql
	
%.run: %.go
	GOPATH=$(GOPATH) $(GO) run $<

deploy:
	rsync http-server http-client server1:
	rsync http-server http-client server2:

clean:
	rm $(TARGETS) version/version.go

help:
	$(GO) help
