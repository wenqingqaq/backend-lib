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
	//code := "d6083ef6ca6e4386edf7"
	//state := "casdoor"
	//token := client.GetOAuthToken(code, state)
	//
	//fmt.Println("------")
	//fmt.Println(token.AccessToken)
	//fmt.Println("------")

	AccessToken := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImNlcnQtYnVpbHQtaW4iLCJ0eXAiOiJKV1QifQ.eyJhdWQiOlsiNmJlMWVjNDk2ZjMxNzNlMzU1MDkiXSwiZW1haWwiOiIxMTY5MTI1MDU1QHFxLmNvbSIsImV4cCI6MTcxMjEzMjE4OCwiaWF0IjoxNzExNTI3Mzg4LCJpZCI6ImM1ZTE0NWMyLTVmNmYtNDlmOC1hZDU2LWExZmRjNTVhYWUxYSIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6ODEwMCIsImp0aSI6ImFkbWluLzM2NDdjZmZhLTVhM2EtNDUxNi04ZDdjLTcxMDk0NGQ2MDBkMCIsIm5hbWUiOiJ3ZW5xaW5nIiwibmJmIjoxNzExNTI3Mzg4LCJub25jZSI6IiIsInBob25lIjoiMTg3MzYwMDcxNTgiLCJzY29wZSI6InJlYWQiLCJzdWIiOiJjNWUxNDVjMi01ZjZmLTQ5ZjgtYWQ1Ni1hMWZkYzU1YWFlMWEiLCJ0YWciOiIiLCJ0b2tlblR5cGUiOiJhY2Nlc3MtdG9rZW4ifQ.IH4bDlnDo9-KsWYCJ-aqzjXkbTlH8hsQbmGm_ggmZlc8ztWbV9xv4e1h3w-r6BmFRJkcSuCE9CcAdcLpsnPLIKBEoIvdVXSg3HXXyToj23uYluUEscdRVtEAR7yCdU2skeUpy-j_BHroPTio_UiBP7KESyNuyx3yFpn8dusErWXy6STit47aSJTfffnhBUhi3hKDg2SYsExwN7IuRdapx90GYcE-biP6waoR1em0bCVjiHPPAPTPHtMoItpBUHS3d9lNeRC62HransamO4-iTMsLHJytboT2e1DCxgDKAF7lsHo4TPnmKS1hkurldyV7Lhz5gF7tMpjWd5X8BSQk4SRRVAt1ADQn2kqjuspM7m25siKi1cWlMW7Hbn59hhVupubjVnt9M9ns5v1tRBGh0ZH3H79KhuJvb2hrwDd7hagfS0u5AoPcmlRSM11w_mTVzaglT5XWIItsxJfZcApOY2C9G71NsvnQ85fZZix-A0gQHjG-YmogT9eHri4pyb_3AIJ7AifDdAoR5nrbx7vniec1zj_xVPjncD-QgCA9qy5nv6Ui-UqLfHXGz065Ikxr0X96TQNtJQnt9oYDCZutjbukpM2bRKP6SKBVqxNPIIy3_J3uVrKJW2XpIXokuKqMmlVS9c7wgy4h_1oCtM2cN9jHW_VLHpgbRho6LuzYzc8"
	// 验证访问令牌
	claims, err := client.Client.ParseJwtToken(AccessToken)
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
