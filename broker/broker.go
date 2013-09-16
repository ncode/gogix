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
	"fmt"
	"github.com/ncode/gogix/syslog"
	"github.com/ncode/gogix/utils"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

type Connection struct {
	conn       *amqp.Connection
	pub        *amqp.Channel
	Queue      string
	Expiration string
	Uri        string
	mu         sync.RWMutex
}

func setup(uri, queue string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		utils.Check(err, "Unable to connect to broker")
		return nil, nil, err
	}

	pub, err := conn.Channel()
	if err != nil {
		utils.Check(err, "Unable to acquire channel")
		return nil, nil, err
	}

	err = pub.ExchangeDeclare(queue, "direct", true, true, false, false, nil)
	if err != nil {
		utils.Check(err, "Unable to declare exchange")
		return nil, nil, err
	}

	_, err = pub.QueueDeclare(queue, true, false, false, false, nil)
	if err != nil {
		utils.Check(err, "Unable to declare queue")
		return nil, nil, err
	}

	err = pub.QueueBind(queue, queue, queue, false, nil)
	if err != nil {
		utils.Check(err, "Unable to bind queue")
		return nil, nil, err
	}

	return conn, pub, nil
}

func (c Connection) SetupBroker() Connection {
	var err error
	c.conn, c.pub, err = setup(c.Uri, c.Queue)
	utils.CheckPanic(err, "Problem acquiring connection")
	return c
}

func (c Connection) Send(parsed syslog.Graylog2Parsed) (err error) {
	encoded, err := json.Marshal(parsed)
	utils.Check(err, "Unable to encode json")
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "text/plain",
		Body:         []byte(encoded),
		Expiration:   c.Expiration,
	}

	//c.mu.Lock()
	//defer c.mu.Unlock()
	fmt.Println("lalala")
	err = c.pub.Publish(c.Queue, c.Queue, false, false, msg)
	if err != nil {
		utils.Check(err, "Unable to publish message")
	}

	return err
}

func (c *Connection) NotifyClose() (err error) {
	bc := make(chan *amqp.Error)
	c.conn.NotifyClose(bc)
	for {
		b := <-bc
		if b != nil {
			for {
				fmt.Println("meh")
				//c.mu.Lock()
				c.conn, c.pub, err = setup(c.Uri, c.Queue)
				if err == nil {
					c.conn.NotifyClose(bc)
					break
				}
				time.Sleep(2 * time.Second)
				//c.mu.Unlock()
			}
		}
	}
}

func (c Connection) Close() {
	defer c.conn.Close()
}
