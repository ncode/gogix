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
	Queue      string
	Expiration string
	Uri        string
}

var (
	max_retries = 40
)

func (c Connection) SetupBroker() (Connection, error) {
	conn, err := amqp.Dial(c.Uri)
	if err != nil {
		utils.Check(err, "Unable to connect to broker")
		return c, err
	}
	c.conn = conn

	pub, err := c.conn.Channel()
	if err != nil {
		utils.Check(err, "Unable to acquire channel")
		return c, err
	}

	c.pub = pub
	err = c.pub.ExchangeDeclare(c.Queue, "direct", true, true, false, false, nil)
	if err != nil {
		utils.Check(err, "Unable to declare queue")
		return c, err
	}

	_, err = c.pub.QueueDeclare(c.Queue, true, false, false, false, nil)
	if err != nil {
		utils.Check(err, "Unable to declare queue")
		return c, err
	}

	err = c.pub.QueueBind(c.Queue, c.Queue, c.Queue, false, nil)
	if err != nil {
		utils.Check(err, "Unable to declare queue")
		return c, err
	}

	return c
}

func (c Connection) Send(parsed syslog.Graylog2Parsed) (err error) {
	encoded, err := json.Marshal(parsed)
	utils.Check(err, "Unable to encode json")
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         []byte(encoded),
		Expiration:   c.Expiration,
	}

	err = c.pub.Publish(c.Queue, c.Queue, false, false, msg)
	if err != nil {
		utils.Check(err, "Unable to publish message")
		time.Sleep(time.Second)
		c = c.SetupBroker()
	}

	return err
}

func (c *Connection) Close() {
	defer c.conn.Close()
}
