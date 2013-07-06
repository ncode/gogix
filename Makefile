.PHONY: build install
DESTDIR=$(CURDIR)/debian/gogix
SRCDIR=$(CURDIR)/src

clean:
	rm --preserve-root -rf $(DESTDIR)

build:
	go get github.com/streadway/amqp
	go get github.com/msbranco/goconfig
	cd $(SRCDIR) && go build -o gogix-server

install:
	install -d --mode=755 $(DESTDIR)
	install -d --mode=755 $(DESTDIR)/usr/sbin
	install -d --mode=755 $(DESTDIR)/etc/gogix
	install -d --mode=755 $(DESTDIR)/etc/rsyslog.d
	install -v --mode=755 $(SRCDIR)/gogix-server $(DESTDIR)/usr/sbin/gogix-server
	install -v --mode=644 $(SRCDIR)/config/gogix.conf $(DESTDIR)/etc/gogix/gogix.conf
	install -v --mode=644 $(SRCDIR)/config/gogix-rsyslog.conf $(DESTDIR)/etc/rsyslog.d/gogix-rsyslog.conf
