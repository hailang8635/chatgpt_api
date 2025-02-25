package main

import (
	"chatgpt_api/gpt_api"
	"chatgpt_api/utils"
	"fmt"
)

func main() {

	//gpt_api.GLM_test()
	//gpt_api.GLM_test()

	utils.TestBase64()

	fmt.Println("start...")
	gpt_api.Gpt_http_server()
	fmt.Println("stop...")
}
