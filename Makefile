.PHONY: build install
DESTDIR=$(CURDIR)/debian/gogix
GOPATH=$(CURDIR)/debian/go-build

clean:
	rm --preserve-root -rf $(DESTDIR)
	rm --preserve-root -rf $(GOPATH)

build:
	GOPATH=$(GOPATH) go get -v github.com/ncode/gogix/...
	GOPATH=$(GOPATH) go build -o gogix

install:
	install -d --mode=755 $(DESTDIR)
	install -d --mode=755 $(DESTDIR)/usr/sbin
	install -d --mode=755 $(DESTDIR)/etc/gogix
	install -d --mode=755 $(DESTDIR)/etc/rsyslog.d
	install -v --mode=755 gogix $(DESTDIR)/usr/sbin/gogix
	install -v --mode=644 config/gogix.conf $(DESTDIR)/etc/gogix/gogix.conf
	install -v --mode=644 config/gogix-rsyslog.conf $(DESTDIR)/etc/rsyslog.d/gogix-rsyslog.conf
