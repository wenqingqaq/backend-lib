package tecmq

import (
	"errors"
	"github.com/nats-io/nats.go"
)

type NatQueueCon struct {
	Url string
	Nc  *nats.Conn
}

func NewNatsSubPushClient(url string) (*NatQueueCon, error) {
	if url == "" {
		return nil, errors.New("nats config is empty ")
	}
	nc, err := NewNatsClient(url)
	if err != nil {
		return nil, err
	}

	return &NatQueueCon{
		Url: url,
		Nc:  nc,
	}, nil
}

func (n *NatQueueCon) CloseNc() {
	n.Nc.Close()
}

func (n *NatQueueCon) Drain() error {
	return n.Nc.Drain()
}

func (n *NatQueueCon) Flush() error {
	return n.Nc.Flush()
}

func (n *NatQueueCon) Push(subject string, b []byte) error {
	msg := &nats.Msg{
		Subject: subject,
		Data:    b,
	}
	err := n.Nc.PublishMsg(msg)
	if err != nil {
		return err
	}
	return err
}

func (n *NatQueueCon) Sub(subject string, f func(msgC *MqMsg)) (*MqSub, error) {
	tecSub := &MqSub{}
	sub, err := n.Nc.Subscribe(subject, HandlerFucAndAck(f)) // nats.ManualAck()
	if err != nil {
		return nil, err
	}
	tecSub.Sub = sub

	return tecSub, nil
}

func (n *NatQueueCon) SubSync(subject string) (*MqSub, error) {
	tecSub := &MqSub{}
	sub, err := n.Nc.SubscribeSync(subject) // nats.ManualAck()
	if err != nil {
		return nil, err
	}
	tecSub.Sub = sub

	return tecSub, nil
}

// QueueSub 订阅一个队列消息
func (n *NatQueueCon) QueueSub(subject, queue string, f func(msgC *MqMsg)) (*MqSub, error) {
	sub, err := n.Nc.QueueSubscribe(subject, queue, HandlerFucAndAck(f)) // nats.ManualAck()
	if err != nil {
		return nil, err
	}

	return &MqSub{Sub: sub}, nil
}

// QueueSubSync 订阅一个js队列消息
func (n *NatQueueCon) QueueSubSync(subject, queue string) (*MqSub, error) {
	sub, err := n.Nc.QueueSubscribeSync(subject, queue)
	if err != nil {
		return nil, err
	}

	return &MqSub{Sub: sub}, nil
}
