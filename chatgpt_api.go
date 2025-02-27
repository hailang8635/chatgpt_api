package main

import (
	"chatgpt_api/wechat_server"
	"log"
)

func main() {

	//gpt_api.GLM_test()
	//gpt_api.GLM_test()

	//utils.TestBase64()

	log.Println("start...")
	wechat_server.Gpt_http_server()
	log.Println("stop...")
}
