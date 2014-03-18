GO=/usr/local/go/bin/go

all: parse-dwml http-server

http-server: http-server.go
	$(GO) build $<

parse-dwml: parse-dwml.go
	$(GO) build $<

%.run: %.go
	$(GO) run $<

clean:
	rm http-server parse-dwml

help:
	$(GO) help
