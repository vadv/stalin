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
	

clean:
	rm -rf $(GOPATH)

install:
	mkdir -p $(DESTDIR)/usr/bin
	cp $(TMPSRC)/$(GOPKG)/$(GOPKG) $(DESTDIR)/usr/bin
	chmod 0755 $(DESTDIR)/usr/bin/$(GOPKG)

.PHONY: all clean install
