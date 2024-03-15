package teckafka

import "github.com/Shopify/sarama"

type IHandlerConsumer interface {
	PreCallBack(msg []byte) error
	CallBack(msg *sarama.ConsumerMessage) error
	PostCallBack(msg []byte) error
	DealErrCallBack(errMsg string)
}

type BaseHandlerConsumer struct {
}

/**
*	调用前接口
 */
func (b *BaseHandlerConsumer) PreCallBack(msg []byte) error {
	return nil
}
func (b *BaseHandlerConsumer) CallBack(msg *sarama.ConsumerMessage) error {
	return nil
}
func (b *BaseHandlerConsumer) PostCallBack(msg []byte) error {
	return nil
}
func (b *BaseHandlerConsumer) DealErrCallBack(errMsg string) {

}

type ConsumerGroupHandler struct {
	Fn IHandlerConsumer
}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	//h.Fn.CallBack内为异步处理，经测试，此处不需要再做异步
	for msg := range claim.Messages() {
		_ = h.Fn.PreCallBack(msg.Value)
		_ = h.Fn.CallBack(msg)
		_ = h.Fn.PostCallBack(msg.Value)
		sess.MarkMessage(msg, "")
	}
	return nil
}
