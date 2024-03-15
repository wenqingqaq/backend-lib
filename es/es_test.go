/*
@author：wj@tecorigin.com
@date：2023/9/13 16:35
@note：
*/
package es

import (
	"fmt"
	"testing"
)

func TestDB(t *testing.T) {
	client, err := NewESClient([]string{"http://10.8.20.1:30200"}, "elastic", "Tc123456")
	fmt.Println("client---", client)
	//fmt.Println("c---", client.)
	get, err := client.API.Indices.Get([]string{"k8s-kubecube-project-2023.09.14"})
	if err != nil {
		fmt.Println("err000---", err)
		return
	}
	s := get.String()
	fmt.Println("aaa--", s)
}
