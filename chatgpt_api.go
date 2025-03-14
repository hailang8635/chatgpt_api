package main

import (
	"chatgpt_api/batch_loop_ask_task"
	"chatgpt_api/config"
	"chatgpt_api/wechat_server"
	"fmt"
	"log"
	"os"
)

func main() {

	//gpt_api.GLM_test()
	//gpt_api.GLM_test()

	//utils.TestBase64()

	config.InitProperties()

	log.Println("start batch...")

	if len(os.Args) >= 2 {
		// 获取所有启动参数（包含程序路径）
		allArgs := os.Args
		fmt.Printf("程序路径: %s\n", allArgs[0])

		// 获取用户参数（排除程序路径）
		userArgs := allArgs[1:]
		fmt.Println("用户参数列表:")
		for i, arg := range userArgs {
			fmt.Printf("参数%d: %s\n\n", i+1, arg)
		}

		//batch_loop_ask_task.Ask_batch("sample_company_1.csv")
		batch_loop_ask_task.Ask_batch(userArgs[0])
	}
	log.Println("end batch...")

	log.Println("start wechat server...")
	// TODO 放开
	//utils.InitGDBC()
	wechat_server.Gpt_http_server()

	log.Println("stop wechat server...")
}
