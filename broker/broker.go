package broker

import (
	"encoding/json"
	"github.com/ncode/gogix/syslog"
	"github.com/ncode/gogix/utils"
	"github.com/streadway/amqp"
	"time"
)

type Connection struct {
	conn       *amqp.Connection
	pub        *amqp.Channel
	queue      string
	expiration string
}

func (self Connection) Dial(uri string) Connection {
	conn, err := amqp.Dial(uri)
	utils.CheckPanic(err, "Unable to connect to broker")
	self.conn = conn
	return self
}

func (self Connection) SetupBroker(queue string, message_ttl string) Connection {
	pub, err := self.conn.Channel()
	utils.CheckPanic(err, "Unable to acquire channel")
	self.pub = pub
	self.expiration = message_ttl
	_, err = self.pub.QueueDeclare(queue, true, false, false, false, nil)
	utils.CheckPanic(err, "Unable to declare queue")
	self.queue = queue
	return self
}

func (self Connection) Send(parsed syslog.Parser) {
	encoded, err := json.Marshal(parsed)
	utils.Check(err, "Unable to encode json")
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         encoded,
		Expiration:   self.expiration,
	}

	err = self.pub.Publish("", self.queue, false, false, msg)
	utils.Check(err, "Unable to publish message")
}

func (self Connection) Close() {
	defer self.conn.Close()
}
