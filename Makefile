.PHONY: build install
DESTDIR=$(CURDIR)/debian/gogix
GOPATH=$(CURDIR)/debian/go-build

clean:
    rm --preserve-root -rf $(DESTDIR)

build:
    GOPATH=$(GOPATH) go get github.com/ncode/gogix
    GOPATH=$(GOPATH) go get github.com/streadway/amqp
    GOPATH=$(GOPATH) go get github.com/msbranco/goconfig
    GOPATH=$(GOPATH) go build -o gogix

install:
    install -d --mode=755 $(DESTDIR)
    install -d --mode=755 $(DESTDIR)/usr/sbin
    install -d --mode=755 $(DESTDIR)/etc/gogix
    install -d --mode=755 $(DESTDIR)/etc/rsyslog.d
    install -v --mode=755 gogix-server $(DESTDIR)/usr/sbin/gogix-server
    install -v --mode=644 config/gogix.conf $(DESTDIR)/etc/gogix/gogix.conf
    install -v --mode=644 config/gogix-rsyslog.conf $(DESTDIR)/etc/rsyslog.d/gogix-rsyslog.conf
