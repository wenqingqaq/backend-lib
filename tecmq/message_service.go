package tecmq

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"time"
)

const (
	ChannelTypeEmail    = 1 // 邮件
	ChannelTypeWebhook  = 2 // WebHook
	ChannelTypeDingDing = 3 // 钉钉

	DefaultMessageSubject = "message-service.message"
)

type Msgs struct {
	Source           uint8      `json:"source"` // 服务来源类型 1-事件监控系统 2-监控告警系统 3-....
	Data             []*MsgData // 发送消息信息
	SourceCreateTime time.Time  // 消息创建时间
}

type MsgData struct {
	Title      string                 `json:"title"`
	Payload    map[string]interface{} `json:"payload"`     // 消息体
	TemplateId int                    `json:"template_id"` // 模板id
	Receivers  []Receiver             `json:"receivers"`   // 消息接收者
	Type       uint8                  `json:"type"`        // 通道类型 1-邮件 2-WebHook
}

type Receiver struct {
	Id      string `json:"contacts_id"`   // 联系人ID
	Name    string `json:"contacts_name"` // 联系人名称
	Address string `json:"address"`       // 通知地址
}

// NewTestMsg 测试消息生成例子
func NewTestMsg() *Msgs {
	return &Msgs{
		Data: []*MsgData{
			{
				Title: "【警告】训练启动失败",
				Payload: map[string]interface{}{
					"NODE":    "tc2",
					"Message": "有个节点报警了",
					"name":    "1",
					"test":    "2",
				},
				TemplateId: 1,
				Type:       ChannelTypeEmail,
				Receivers: []Receiver{
					{
						Id:      "1001",
						Name:    "1001_name",
						Address: "yanwenqing@163.com",
					},
				},
			},
		},
		Source:           1,
		SourceCreateTime: time.Now(),
	}
}

func HandlerFucAndAck(f func(msgC *MqMsg)) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		f(&MqMsg{Msg: msg})
	}
}

// JsHandlerFucAndAck func(msg Msg)
func JsHandlerFucAndAck(f func(msgC *JsMqMsg)) func(msg jetstream.Msg) {
	return func(msg jetstream.Msg) {
		f(&JsMqMsg{Msg: msg})
	}
}

// PushMessage to MessageService 发送消息
func (n *NatQueueCon) PushMessage(subject string, msgD *Msgs) error {
	b, err := json.Marshal(msgD)
	if err != nil {
		return err
	}
	if subject == "" {
		subject = DefaultMessageSubject
	}
	msg := &nats.Msg{
		Subject: subject,
		Data:    b,
	}
	err = n.Nc.PublishMsg(msg)
	if err != nil {
		return err
	}
	return err
}
