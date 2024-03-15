package teckafka

import (
	"context"
	"gitee.com/yanwenqing/backend-lib/logz"
	"github.com/Shopify/sarama"
)

/**
*	@param  consumerGroupHandler 结构体ConsumerGroupHandler的子类
*	@param  client kafka集群
*	@param  group 消费者组
*	@param 	topic 消费者topic
 */
func StartConsumer(consumerGroupHandler ConsumerGroupHandler, version sarama.KafkaVersion, kafkaHosts []string, saslEnable bool, userName, password, group, topic string) {
	config := sarama.NewConfig()
	config.Version = version
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest // 如果以前没有提交偏移量，则要使用的偏移量是最旧偏移量。
	if true == saslEnable {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = userName
		config.Net.SASL.Password = password
	}
	//config.Consumer.Fetch.Min = 5
	//config.Consumer.Fetch.Max = 10
	// Start with a client
	client, err := sarama.NewClient(kafkaHosts, config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = client.Close() }()
	// Start a new consumer group
	g, err := sarama.NewConsumerGroupFromClient(group, client)
	if err != nil {
		panic(err)
	}
	defer func() { _ = g.Close() }()

	go func() {
		for err := range g.Errors() {
			if err != nil {
				consumerGroupHandler.Fn.DealErrCallBack(err.Error())
			} else {
				logz.Error("kafka err is nil ...")
			}
		}
	}()
	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		topics := []string{topic}
		handler := consumerGroupHandler
		err := g.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}
	}
}
