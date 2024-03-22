package tecmq

import (
	"encoding/json"
	"time"
)

const (
	SourceDigital         = 1
	SourceDigitalProduce  = 2
	DigitalProduceSubject = "digital.produce"
)

type DigitalProduceMsg struct {
	Source           uint8                    `json:"source"` // 服务来源类型 1-事件监控系统 2-监控告警系统 3-....
	Data             []*DigitalProduceMsgData // 发送消息信息
	SourceCreateTime time.Time                // 消息创建时间
}

// DigitalProduceMsgData 订阅主题使用的消息
type DigitalProduceMsgData struct {
	BackgroundUri       string `json:"background_uri"`
	FaceVideo           string `json:"face_video"`
	AudionUri           string `json:"audion_uri"`
	StoryboardId        int32  `json:"storyboard_id"`
	StoryboardSectionId int32  `json:"storyboard_section_id"`
	Payload             string `json:"payload"`
}

type ProduceResultMsg struct {
	Source           uint8               `json:"source"` // 服务来源类型 1-事件监控系统 2-监控告警系统 3-....
	Data             *ProduceCallbackMsg // 发送消息信息
	SourceCreateTime time.Time           // 消息创建时间
}
type ProduceCallbackMsg struct {
	StoryboardId int32
	OutFile      string
	IsSuccess    int32
	Duration     string
}

// PushProduceMessage to ProduceServer 发送生产视频的消息
func (n *NatJs) PushProduceMessage(subject string, msgD *DigitalProduceMsg) error {
	b, err := json.Marshal(msgD)
	if err != nil {
		return err
	}
	err = n.JsPushMessage(subject, b)
	if err != nil {
		return err
	}

	return err
}

func (n *NatJs) PushProduceCallBack(subject string, msgD *ProduceResultMsg) error {
	b, err := json.Marshal(msgD)
	if err != nil {
		return err
	}
	err = n.JsPushMessage(subject, b)
	if err != nil {
		return err
	}

	return err
}
