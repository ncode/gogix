.PHONY: build install
DESTDIR=$(CURDIR)/debian/gogix

clean:
	rm --preserve-root -rf $(DESTDIR)

build:
	go get github.com/ncode/gogix
	go get github.com/streadway/amqp
	go get github.com/msbranco/goconfig
	go build -o gogix

install:
	install -d --mode=755 $(DESTDIR)
	install -d --mode=755 $(DESTDIR)/usr/sbin
	install -d --mode=755 $(DESTDIR)/etc/gogix
	install -d --mode=755 $(DESTDIR)/etc/rsyslog.d
	install -v --mode=755 gogix-server $(DESTDIR)/usr/sbin/gogix-server
	install -v --mode=644 config/gogix.conf $(DESTDIR)/etc/gogix/gogix.conf
	install -v --mode=644 config/gogix-rsyslog.conf $(DESTDIR)/etc/rsyslog.d/gogix-rsyslog.conf
