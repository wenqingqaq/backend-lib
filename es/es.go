/*
@author：wj@tecorigin.com
@date：2023/9/13 15:38
@note：ES 客户端
*/
package es

import (
	"gitee.com/yanwenqing/backend-lib/logz"
	"github.com/elastic/go-elasticsearch/v7"
	"go.uber.org/zap"
)

func NewESClient(address []string, username, password string) (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		Addresses: address,
		Username:  username,
		Password:  password,
	}
	es, err := elasticsearch.NewClient(cfg)
	//fmt.Println("err---", err)
	if err != nil {
		logz.Fatal("创建ES客户端失败：", zap.Any("error：", err))
	}
	return es, err
}
