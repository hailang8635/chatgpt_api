package gpt_api

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

/**
 * 命令行使用
 */
func Ask_gpt() (string, error) {
	if len(os.Args) >= 2 {
		args1 := os.Args[1]
		respBody, err := GptApi(args1)

		if err != nil {
			log.Println("gpt_api 调用失败")
			return "", err
		}
		log.Println(respBody)
		log.Println("请输入您要问的内容？【N】退出")

	} else {
		log.Println("请输入您要问的内容？")
	}
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()

		// 遇到N退出
		if line == "N" {
			break
		}

		log.Println("正在问ChatGpt ... 【", line, "】")
		resp, err := GptApi(line)
		if err != nil {
			log.Println("gpt_api 调用失败")
			return "", err
		} else {
			log.Println(resp)
		}
		//log.Println()
		fmt.Println()
		log.Println("请继续输入您要问的内容？")
		fmt.Println()
	}

	return "ask done", nil
}
