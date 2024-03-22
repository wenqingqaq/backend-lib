package tecmq

import (
	"errors"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"time"
)

type MqMsg struct {
	Msg *nats.Msg
}

type MsgHandlerFunc func(msg *nats.Msg) // 回调处理方法

type JsMqMsg struct {
	Msg jetstream.Msg
}

type JsMsgHandlerFunc func(msg jetstream.Msg) // 回调处理方法

type MqSub struct {
	Sub *nats.Subscription
}

func NatsJsClient(url string, jsName string, subjects []string) (*nats.Conn, nats.JetStreamContext, error) {
	nc, err := NewNatsClient(url)
	if err != nil {
		return nil, nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, nil, err
	}

	s, err := js.StreamInfo(jsName)
	if s == nil && jsName != "" && len(subjects) > 0 { // 没有js 新增加一个
		cfg := nats.StreamConfig{
			Name:     jsName,
			Subjects: subjects,
			MaxMsgs:  10000, // 单个主题最大消息1w
		}
		cfg.Storage = nats.FileStorage
		_, err = js.AddStream(&cfg)
		if err != nil {
			return nil, nil, err
		}
		return nc, js, nil
	}

	return nc, js, nil
}

func NewNatsClient(url string) (*nats.Conn, error) {
	if url == "" {
		url = nats.DefaultURL
	}
	opts := nats.Options{
		Url:            url,
		AllowReconnect: true,
		MaxReconnect:   100,
		ReconnectWait:  5 * time.Second,
		Timeout:        1 * time.Second,
	}

	nc, err := opts.Connect()
	if err != nil {
		return nil, err
	}

	return nc, nil
}

func NewEncodedNatsClient(url string) (*nats.EncodedConn, error) {
	if url == "" {
		return nil, errors.New("url can't be empty")
	}
	opts := nats.Options{
		Url:            url,
		AllowReconnect: true,
		MaxReconnect:   100,
		ReconnectWait:  5 * time.Second,
		Timeout:        1 * time.Second,
	}

	nc, err := opts.Connect()
	if err != nil {
		return nil, err
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	return ec, nil
}
