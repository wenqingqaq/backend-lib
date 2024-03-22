package tecmq

import (
	"context"
	"errors"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type NatJs struct {
	Url     string
	JsName  string
	Subject []string
	JsCtx   nats.JetStreamContext
	Js      jetstream.JetStream
	Nc      *nats.Conn
}

func NewNatsJsSubPushQueueClient(url, jsName string, subject []string) (*NatJs, error) {
	if url == "" || jsName == "" || len(subject) == 0 {
		return nil, errors.New("nats config is empty ")
	}
	nc, jsc, err := NatsJsClient(url, jsName, subject)
	if err != nil {
		return nil, err
	}
	if jsc == nil {
		return nil, errors.New("js is empty , has no conn ")
	}
	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	return &NatJs{
		Url:     url,
		JsName:  jsName,
		Subject: subject,
		Nc:      nc,
		JsCtx:   jsc,
		Js:      js,
	}, nil
}

func (n *NatJs) CloseNc() {
	n.Nc.Close()
}

func (n *NatJs) Drain() {
	if n.Nc != nil {
		err := n.Nc.Drain()
		if err != nil {
			return
		}
	}
}

func (n *NatJs) Flush() {
	if n.Nc != nil {
		err := n.Nc.Flush()
		if err != nil {
			return
		}
	}
}

func (n *NatJs) CreateOrUpdateConsumer(conName string) (jetstream.Consumer, error) {
	if n.JsName == "" {
		return nil, errors.New("has no js name")
	}

	con, err := n.Js.CreateOrUpdateConsumer(context.Background(), n.JsName, jetstream.ConsumerConfig{
		Durable:   conName,
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return nil, err
	}

	return con, nil
}

// JsPushMessage 发送消息
func (n *NatJs) JsPushMessage(subject string, b []byte) error {
	if n.Js == nil {
		return errors.New("js conn is nil")
	}
	msg := &nats.Msg{
		Subject: subject,
		Data:    b,
	}
	_, err := n.Js.PublishMsg(context.Background(), msg)
	if err != nil {
		return err
	}
	return err
}

func (n *NatJs) JsConsumerSub(conName string, f func(msg *JsMqMsg)) (jetstream.ConsumeContext, error) {
	consumer, err := n.Js.Consumer(context.Background(), n.JsName, conName)
	if err != nil {
		return nil, err
	}
	c, err := consumer.Consume(JsHandlerFucAndAck(f), jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
		fmt.Println(err)
	}))
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (n *NatJs) JsConsumerSubStop(c jetstream.ConsumeContext) {
	c.Stop()
}

// JsQueueSub 订阅一个js队列消息
func (n *NatJs) JsQueueSub(subject, queue, consumer string, f func(msgC *MqMsg)) (*MqSub, error) {
	if n.Js == nil {
		return nil, errors.New("js conn is nil")
	}

	// MaxDeliver 消息最大投递次数目前没有限制；
	// nats.ManualAck() 消息需要确切的ack才认为是正确消费的 nats.MaxDeliver(3)
	sub, err := n.JsCtx.QueueSubscribe(subject, queue, HandlerFucAndAck(f),
		nats.Durable(consumer), nats.ManualAck(),
	)
	if err != nil {
		return nil, err
	}
	return &MqSub{Sub: sub}, nil
}

func (n *NatJs) JsQueueSubSync(subject, queue string) (*MqSub, error) {
	if n.JsCtx == nil {
		return nil, errors.New("js conn is nil")
	}

	sub, err := n.JsCtx.QueueSubscribeSync(subject, queue, nats.OrderedConsumer())
	if err != nil {
		return nil, err
	}
	return &MqSub{Sub: sub}, nil
}
