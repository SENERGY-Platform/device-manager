/*
 * Copyright 2019 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package listener

import (
	"context"
	"errors"
	"github.com/SENERGY-Platform/device-manager/lib/kafka/util"
	"github.com/segmentio/kafka-go"
	"io"
	"log"
	"sync"
	"time"
)

func NewConsumer(ctx context.Context, broker string, groupid string, topic string, listener func(topic string, msg []byte) error, errorhandler func(err error, consumer *Consumer)) (consumer *Consumer, err error) {
	consumer = &Consumer{ctx: ctx, groupId: groupid, broker: broker, topic: topic, listener: listener, errorhandler: errorhandler}
	err = consumer.start()
	return
}

type Consumer struct {
	count        int
	broker       string
	groupId      string
	topic        string
	ctx          context.Context
	listener     func(topic string, msg []byte) error
	errorhandler func(err error, consumer *Consumer)
	mux          sync.Mutex
}

func (this *Consumer) start() error {
	log.Println("DEBUG: consume topic: \"" + this.topic + "\"")

	err := util.InitTopic(this.broker, this.topic)
	if err != nil {
		log.Println("ERROR: unable to create topic", err)
		return err
	}
	r := kafka.NewReader(kafka.ReaderConfig{
		CommitInterval:         0, //synchronous commits
		Brokers:                []string{this.broker},
		GroupID:                this.groupId,
		Topic:                  this.topic,
		MaxWait:                1 * time.Second,
		Logger:                 log.New(io.Discard, "", 0),
		ErrorLogger:            log.New(io.Discard, "", 0),
		WatchPartitionChanges:  true,
		PartitionWatchInterval: time.Minute,
	})
	go func() {
		defer r.Close()
		defer log.Println("close consumer for topic ", this.topic)
		for {
			select {
			case <-this.ctx.Done():
				return
			default:
				m, err := r.FetchMessage(this.ctx)
				if err == io.EOF || err == context.Canceled {
					return
				}
				if err != nil {
					log.Println("ERROR: while consuming topic ", this.topic, err)
					this.errorhandler(err, this)
					return
				}

				err = retry(func() error {
					return this.listener(m.Topic, m.Value)
				}, func(n int64) time.Duration {
					return time.Duration(n) * time.Second
				}, 10*time.Minute)

				if err != nil {
					log.Println("ERROR: unable to handle message (no commit)", err)
					this.errorhandler(err, this)
				} else {
					err = r.CommitMessages(this.ctx, m)
				}
			}
		}
	}()
	return err
}

func retry(f func() error, waitProvider func(n int64) time.Duration, timeout time.Duration) (err error) {
	err = errors.New("")
	start := time.Now()
	for i := int64(1); err != nil && time.Since(start) < timeout; i++ {
		err = f()
		if err != nil {
			log.Println("ERROR: kafka listener error:", err)
			wait := waitProvider(i)
			if time.Since(start)+wait < timeout {
				log.Println("ERROR: retry after:", wait.String())
				time.Sleep(wait)
			} else {
				return err
			}
		}
	}
	return err
}
