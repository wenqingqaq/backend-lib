package teckafka

import (
	"fmt"
	"gitee.com/yanwenqing/backend-lib/logz"
	"github.com/Shopify/sarama"
	"sync"
)

/**
*	单例模式获取
 */
var (
	defaultKafkaOnce          sync.Once
	defaultKafkaAsyncProducer sarama.AsyncProducer
)

type ErrFn func(err *sarama.ProducerError)
type SuccessFn func(msg *sarama.ProducerMessage)

/**
*	默认回调错误处理方式
 */
func DefaultErrFn(err *sarama.ProducerError) {
	if err != nil {
		logz.Error(fmt.Sprintf("kafka err:%+v", err.Error()))
	}

}

/**
*	默认回调成功处理方式
 */
func DefaultSuccessFn(msg *sarama.ProducerMessage) {
	logz.Info(fmt.Sprintf("%+v", msg))
}

/**
*	获取kafka异步生产者
*	@param:connType 获取实例类型
 */
func kafkaAsyncProducerFactory(host []string, sFn SuccessFn, eFn ErrFn) (asyncProducer sarama.AsyncProducer, err error) {
	return getDefaultKafkaProducer(host, false, "", "", sFn, eFn)
}

/**
*	获取默认的kafka实例
*	@param:address []string {} 集群地址
 */
func getDefaultKafkaProducer(address []string, saslEnable bool, userName, password string, sFn SuccessFn, eFn ErrFn) (asyncProducer sarama.AsyncProducer, err error) {
	defaultKafkaOnce.Do(func() {
		config := sarama.NewConfig()
		config.Producer.RequiredAcks = sarama.WaitForAll
		config.Producer.Partitioner = sarama.NewRandomPartitioner
		if true == saslEnable {
			config.Net.SASL.Enable = true
			config.Net.SASL.User = userName
			config.Net.SASL.Password = password
		}
		config.Producer.Return.Successes = true
		config.Producer.Return.Errors = true
		defaultKafkaAsyncProducer, err = sarama.NewAsyncProducer(address, config)
		if err != nil {
			logz.Error(fmt.Sprintf("getDefaultKafkaProducer err:%+v", err))
		}
		go startListen(defaultKafkaAsyncProducer, sFn, eFn)
	})
	return defaultKafkaAsyncProducer, err
}

/**
*	发送消息
*	@param:topic string 所属topic
*	@param:msg string 消息
*	@param:key sarama.Encoder
 */
func Publish(host []string, topic, msg string, key sarama.Encoder, s SuccessFn, e ErrFn) (sourceMsg string, err error) {
	//获取kafka实例
	producerAsync, err := kafkaAsyncProducerFactory(host, s, e)
	if err != nil {
		return msg, err
	}
	kMsg := &sarama.ProducerMessage{Topic: topic, Key: key, Value: sarama.StringEncoder(msg)}
	//使用通道发送消息
	producerAsync.Input() <- kMsg
	return msg, nil
}

func startListen(p sarama.AsyncProducer, sFn SuccessFn, eFn ErrFn) {
	for {
		select {
		case suc := <-p.Successes():
			sFn(suc)
		case fail := <-p.Errors():
			eFn(fail)
		}
	}
}
