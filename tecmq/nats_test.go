package tecmq

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNatCon(t *testing.T) {
	nc, err := NewNatsClient("nats://10.8.20.1:30866")
	assert.Empty(t, err)
	s := nc.Status()
	assert.NotEmpty(t, s)
	fmt.Println(nc.Status())

	js, err := nc.JetStream()
	assert.Empty(t, err)
	fmt.Println(js.StreamInfo("MESSAGE3"))
}

func TestNcPush(t *testing.T) {
	nc, err := NewNatsSubPushClient("nats://10.8.20.1:30866")
	assert.Empty(t, err)
	defer nc.CloseNc()
	err = nc.Push("test.log", []byte("test 123"))
	assert.Empty(t, err)
}

func TestNcSub(t *testing.T) {
	nc, err := NewNatsSubPushClient("nats://10.8.20.1:30866")
	assert.Empty(t, err)
	defer nc.CloseNc()
	_, err = nc.Sub("test.log", func(msgC *MqMsg) {
		fmt.Println(string(msgC.Msg.Data))

		msgC.Msg.Ack()
	})
	assert.Empty(t, err)
	select {}
}

func TestNcSubSync(t *testing.T) {
	nc, err := NewNatsSubPushClient("nats://10.8.20.1:30866")
	assert.Empty(t, err)
	defer nc.CloseNc()
	sub, err := nc.SubSync("test.log")
	assert.Empty(t, err)
	msg, err := sub.Sub.NextMsg(time.Second * 10)
	assert.Empty(t, err)
	fmt.Println(string(msg.Data))
	select {}
}

func TestNcSubQueue(t *testing.T) {
	nc, err := NewNatsSubPushClient("nats://10.8.20.1:30866")
	assert.Empty(t, err)
	defer nc.CloseNc()
	_, err = nc.QueueSub("test.log", "test-queue", func(msgC *MqMsg) {
		fmt.Println(string(msgC.Msg.Data))
		msgC.Msg.Ack()
	})
	assert.Empty(t, err)
	select {}
}

func TestJsCon(t *testing.T) {
	nj, err := NewNatsJsSubPushQueueClient("nats://10.8.20.1:30866", "TEST", []string{"test"})
	assert.Empty(t, err)
	s, err := nj.JsCtx.StreamInfo("TEST")
	assert.Empty(t, err)
	assert.NotEmpty(t, s)
}

// 持久化消费
func TestJsSub(t *testing.T) {
	nj, err := NewNatsJsSubPushQueueClient("nats://10.8.20.1:30866", "TEST", []string{"test.*"})
	assert.Empty(t, err)
	_, err = nj.JsQueueSub("test-js", "test-js-queue", "con1", func(msgC *MqMsg) {
		fmt.Println(string(msgC.Msg.Data))
		msgC.Msg.Ack()
	})
	assert.Empty(t, err)
	select {}
}

func TestJsPush(t *testing.T) {
	nj, err := NewNatsJsSubPushQueueClient("nats://10.8.20.1:30866", "TEST", []string{"test.*"})
	assert.Empty(t, err)
	msg := NewTestMsg()
	b, err := json.Marshal(msg)
	assert.Error(t, err)
	err = nj.JsPushMessage("test-js", b)
	assert.Empty(t, err)
}
