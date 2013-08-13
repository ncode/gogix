package broker

import (
	"encoding/json"
	"github.com/ncode/gogix/syslog"
	"github.com/ncode/gogix/util"
	"github.com/streadway/amqp"
	"time"
)

type Connection struct {
	conn       *amqp.Connection
	pub        *amqp.Channel
	queue      string
	expiration int32
}

func (self Connection) Dial(uri string) Connection {
	conn, err := amqp.Dial(uri)
	util.CheckPanic(err, "Unable to connect to broker")
	self.conn = conn
	return self
}

func (self Connection) SetupBroker(queue string, message_ttl int64) Connection {
	pub, err := self.conn.Channel()
	util.CheckPanic(err, "Unable to acquire channel")
	self.pub = pub
	self.expiration = int32(message_ttl)
	opts := amqp.Table{}
	_, err = self.pub.QueueDeclare(queue, true, false, false, false, opts)
	util.CheckPanic(err, "Unable to declare queue")
	self.queue = queue
	return self
}

func (self Connection) Send(parsed syslog.Parser) {
	encoded, err := json.Marshal(parsed)
	util.CheckPanic(err, "Unable to encode json")
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         encoded,
		Expiration:   self.expiration
	}

	err = self.pub.Publish("", self.queue, false, false, msg)
	util.CheckPanic(err, "Unable to publish message")
	defer self.conn.Close()
}
