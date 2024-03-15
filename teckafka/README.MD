1. kafka类库
2. 消费者使用方式
    1. 继承结构体
        ```go
       type SelfStruct struct {
            BaseHandlerConsumer       
       }
       ```
    2. 重写结构体方法
        ```go
        //处理业务前
        func (b *SelfStruct) PreCallBack(msg []byte) error {
            return nil
        }
        //处理业务
        func (b *SelfStruct) CallBack(msg []byte) error {
            return nil
        }
        //处理业务后
        func (b *SelfStruct) PostCallBack(msg []byte) error {
            return nil
        }
        //错误处理
        func (b *SelfStruct) DealErrCallBack(errMsg string) {
        
        }
        ```
    3. 启动消费服务
        ```go
        StartConsumer(SelfStruct{}, gkafka.V1_0_0_0,  []string{"ip1","ip2"}, "groupName", "topicName")
        ```
3. kafka生产者调用方式
    1. 自定义两个方法、分别是处理正确回调和错误回调函数、如果不自定义则调用两个默认处理函数
        ```go
            /**
            *	自定义错误回调处理方式
             */
            func errFn(err *sarama.ProducerError) {
                logger.Error(fmt.Sprintf("err:++++%+v++++", err))
            }
           /**
           *	自定义成功回调处理方式
            */
           func successFn(msg *sarama.ProducerMessage) {
            logger.Info(fmt.Sprintf("success:----%+v------", msg))
           }
        ```
   
     2. 发送消息
        1. 使用自定义函数
            ```go
            参数: 实例  topic  消息体 key 成功回调函数 失败回调函数
            msg, err := Publish(0, "test", "test", nil, successFn, errFn)
            if err != nil {
               fmt.Println(err)
               return
            }
            ```
        2. 使用默认函数
            ```go
            参数: 实例  topic  消息体 key 成功回调函数 失败回调函数
            _, err := Publish(0, "test", "test", nil, DefaultSuccessFn, DefaultErrFn)
            ```
             
           
       
         
         