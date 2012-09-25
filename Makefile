.PHONY: build install
DESTDIR=$(CURDIR)/debian/gogix
SRCDIR=$(CURDIR)/src

clean:
	rm --preserve-root -rf $(DESTDIR)

build:
	go get github.com/streadway/amqp
	go get github.com/kless/goconfig/config
	cd $(SRCDIR) && go build -o logix-server

install:
	install -d --mode=755 $(DESTDIR)
	install -d --mode=755 $(DESTDIR)/usr/sbin
	install -d --mode=755 $(DESTDIR)/etc/logix
	install -d --mode=755 $(DESTDIR)/etc/rsyslog.d
	install -v --mode=755 $(SRCDIR)/logix-server $(DESTDIR)/usr/sbin/logix-server
	install -v --mode=644 $(SRCDIR)/config/logix.conf $(DESTDIR)/etc/logix/logix.conf
	install -v --mode=644 $(SRCDIR)/config/logix-rsyslog.conf $(DESTDIR)/etc/rsyslog.d/logix-rsyslog.conf
