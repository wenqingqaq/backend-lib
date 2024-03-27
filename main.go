package main

import (
	"fmt"
	"github.com/wenqingqaq/backend-lib/casdoor_cx"
)

func main() {
	client, err := casdoor_cx.NewCasDoorClient()
	if err != nil {
		panic(err)
	}
	code := "2e3173c89fced828cd92"
	state := "casdoor"
	token := client.GetOAuthToken(code, state)

	// 验证访问令牌
	claims, err := client.Client.ParseJwtToken(token.AccessToken)
	if err != nil {
		panic(err)
	}
	user, err := client.Client.GetUser(claims.Name)
	if err != nil {
		panic(err)
	}
	fmt.Println(claims.Name)
	fmt.Println(user.Name)
	fmt.Println(user.Permissions)
	for _, item := range user.Permissions {
		fmt.Println(item)
	}
}
