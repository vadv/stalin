export GOPATH := $(CURDIR)/.build
GOPKG := stalin
TMPSRC := $(GOPATH)/src

all: build

build: clean
	mkdir -p $(GOPATH)/src $(GOPATH)/bin $(GOPATH)/pkg
	mkdir -p $(TMPSRC)
	cp -av $(GOPKG) $(TMPSRC)
	cp -av dependencies/code.google.com $(TMPSRC)
	cp -av dependencies/github.com $(TMPSRC)
	cd $(TMPSRC)/$(GOPKG) && GOPATH=$(GOPATH) go build -x -v $(GOPKG)

run: build
	GOGCTRACE=1 $(TMPSRC)/$(GOPKG)/$(GOPKG) -config=$(CURDIR)/example/example.conf

clean:
	rm -rf $(GOPATH)

install:
	mkdir -p $(DESTDIR)/usr/bin
	cp $(TMPSRC)/$(GOPKG)/$(GOPKG) $(DESTDIR)/usr/bin
	chmod 0755 $(DESTDIR)/usr/bin/stalin-client

.PHONY: all clean install
