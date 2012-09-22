package broker

import (
    "time"
    "encoding/json"
    "../../logix/util"
    "../../logix/syslog"
    "github.com/streadway/amqp"
)

type Connection struct {
    conn *amqp.Connection
    pub *amqp.Channel
    queue string
}

func (self Connection) Dial(uri string) Connection{
    conn, err := amqp.Dial(uri)
    util.Checkp(err)
    self.conn = conn
    return self
}

func (self Connection) SetupBroker(queue string) Connection{
    pub, err := self.conn.Channel()
    util.Checkp(err)
    self.pub = pub
    _, err = self.pub.QueueDeclare(queue, true, false, false, false, nil)
    util.Checkp(err)
    self.queue = queue
    return self
}

func (self Connection) Send(parsed syslog.Parser) {
    encoded, err := json.Marshal(parsed)
    util.Checkp(err)
    msg := amqp.Publishing{
        DeliveryMode: amqp.Persistent,
        Timestamp:    time.Now(),
        ContentType:  "text/plain",
        Body:         encoded,
    }

    err = self.pub.Publish("", self.queue, false, false, msg)
    util.Checkp(err)
    defer self.conn.Close()
}