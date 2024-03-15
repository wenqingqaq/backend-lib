package teckafka

import (
	"fmt"
	"gitee.com/yanwenqing/backend-lib/logz"
	"github.com/Shopify/sarama"
	"testing"
	"time"
)

/**
*	自定义错误回调处理方式
 */
func errFn(err *sarama.ProducerError) {
	logz.Error(fmt.Sprintf("err:++++%+v++++", err))
}

/**
*	自定义成功回调处理方式
 */
func successFn(msg *sarama.ProducerMessage) {
	logz.Info(fmt.Sprintf("success:----%+v------", msg))
}

func TestPublish(t *testing.T) {
	msg, err := Publish([]string{"10.8.20.2:32615"}, "push_test", "test11221", nil, successFn, errFn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(msg)
	time.Sleep(time.Second * 1000)
}

func BenchmarkPublish(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Publish([]string{"10.8.20.2:32615"}, "test", "test", nil, DefaultSuccessFn, DefaultErrFn)
		if err != nil {
			fmt.Println(err)
			return
		}

	}
}
