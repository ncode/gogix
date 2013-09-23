# gogix - Transports your local syslog to Graylog2 via AMQP
* gogix is a <a href="https://github.com/ncode/logix">logix</a> port to Go
* http://graylog2.org/about

## Why should I use it?

When you are sending lots of udp log events over the network packet loss can happen, or even
using a tcp log sender you can get a slow response on your server depending on how much logs
your remote log server is receiving simultaneously.

So... what can you do to avoid it?

gogix can help you using its daemon receiving your log events
and queueing you messages on AMQP. You can easily get rid of log event
losses caused by udp and any performance issue that could be caused by
concurrency using tcp remote syslog.

gogix queues your log events on any AMQP Server and you can easy setup
your <a href="https://github.com/Graylog2/graylog2-server">graylog2-server</a> to consume this queue and index your logs on demand.

## Usage:
### Setup your AMQP and Graylog2
* http://www.rabbitmq.com/getstarted.html
* https://github.com/Graylog2/graylog2-server/wiki/AMQP

### Add to your grailog2.conf

    # AMQP
    amqp_enabled = true
    amqp_subscribed_queues = gogix:gelf
    amqp_host = localhost
    amqp_port = 5672
    amqp_username = guest
    amqp_password = guest
    amqp_virtualhost = /

### gogix.conf

    [transport]
    url = amqp://127.0.0.1:5672
    queue = gogix

    [server]
    bind_addr = 127.0.0.1:6660

### on MacOS X:

    $ vim /etc/syslog.conf
    *.notice;authpriv,remoteauth,ftp,install,internal.none  @127.0.0.1:6660
    $ launchctl unload /System/Library/LaunchDaemons/com.apple.syslogd.plist
    $ launchctl load /System/Library/LaunchDaemons/com.apple.syslogd.plist

### on Linux:

    $ vim /etc/rsyslog.d/gogix.conf
    *.*  @127.0.0.1:6660
    $ /etc/init.d/rsyslog restart

### Building

    $ git clone git@github.com:ncode/gogix.git
    $ cd gogix
    $ go get -v .
    $ go build -o gogix

### Building package on Debian and Ubuntu

    $ apt-get install golang
    $ git clone git@github.com:ncode/gogix.git
    $ cd gogix
    $ dpkg-buildpackage -us -uc -rfakeroot

### Running:

    $ Usage: gogix
    $   -h help
    $   -u username
    $   -d debug

    $ GOGIX_CONF=config/gogix.conf ./gogix-server -u $USER -d &
    $ logger test

## Depends:
* amqp - https://github.com/streadway/amqp
* gocofig - https://github.com/msbranco/goconfig
