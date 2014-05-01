GOPATH=$(HOME)/src/gocode
#GO=/usr/local/go/bin/go
GO=$(HOME)/src/go/bin/go
GOROOT=$(HOME)/src/go

TARGETS=\
	parse-dwml \
	http-server http-client \
	ftp-client \
	mysql-client \
	parse-json \
	cache-test \
	tcp-server \
	server-c

all: $(TARGETS)

version/version.go: *.go Makefile
	@./gen-version.sh

%: %.go
	GOROOT=$(GOROOT) GOPATH=$(GOPATH) $(GO) build $<

server-c: server.c
	gcc -g -o server-c server.c -lrt

http-server: version/version.go http-server.go http-server-config.go daemon.go
	GOROOT=$(GOROOT) GOPATH=$(GOPATH) $(GO) build http-server.go http-server-config.go daemon.go

tcp-server: tcp-server.go
	GOROOT=$(GOROOT) GOPATH=$(GOPATH) $(GO) build tcp-server.go

http-client: version/version.go http-client.go http-client-config.go
	GOROOT=$(GOROOT) GOPATH=$(GOPATH) $(GO) build http-client.go http-client-config.go

getdeps:
	GOROOT=$(GOROOT) GOPATH=$(GOPATH) $(GO) get github.com/jlaffaye/ftp
	GOROOT=$(GOROOT) GOPATH=$(GOPATH) $(GO) get github.com/go-sql-driver/mysql
	
%.run: %.go
	GOROOT=$(GOROOT) GOPATH=$(GOPATH) $(GO) run $<

deploy: http-server server-c
	rsync -vz server-c http-server http-server.conf server1:
	rsync -vz http-client server2:

clean:
	rm $(TARGETS) version/version.go

help:
	$(GO) help
