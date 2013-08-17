/* Copyright 2013 Juliano Martinez
   All Rights Reserved.

     Licensed under the Apache License, Version 2.0 (the "License");
     you may not use this file except in compliance with the License.
     You may obtain a copy of the License at

         http://www.apache.org/licenses/LICENSE-2.0

     Unless required by applicable law or agreed to in writing, software
     distributed under the License is distributed on an "AS IS" BASIS,
     WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
     See the License for the specific language governing permissions and
     limitations under the License.

   @author: Juliano Martinez */

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
	self.expiration = message_ttl
	pub, err := self.conn.Channel()
	utils.CheckPanic(err, "Unable to acquire channel")
	self.pub = pub
	err = self.pub.ExchangeDeclare(queue, "direct", true, true, false, false, nil)
	utils.CheckPanic(err, "Unable to declare queue")
	_, err = self.pub.QueueDeclare(queue, true, false, false, false, nil)
	utils.CheckPanic(err, "Unable to declare queue")
	err = self.pub.QueueBind(queue, queue, queue, false, nil)
	utils.CheckPanic(err, "Unable to declare queue")

	self.queue = queue
	return self
}

func (self Connection) Send(parsed syslog.Graylog2Parsed) {
	encoded, err := json.Marshal(parsed)
	utils.Check(err, "Unable to encode json")
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         []byte(encoded),
		Expiration:   self.expiration,
	}

	err = self.pub.Publish(self.queue, self.queue, false, false, msg)
	utils.Check(err, "Unable to publish message")
}

func (self Connection) Close() {
	defer self.conn.Close()
}
