package main

import (
	"chatgpt_api/wechat_server"
	"fmt"
)

func main() {

	//gpt_api.GLM_test()
	//gpt_api.GLM_test()

	//utils.TestBase64()

	fmt.Println("start...")
	wechat_server.Gpt_http_server()
	fmt.Println("stop...")
}
