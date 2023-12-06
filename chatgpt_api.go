package main

import (
    "chatgpt_api/gpt_api"
    "chatgpt_api/utils"
    "fmt"
)

func main() {

    utils.TestBase64()

    fmt.Println("start...")
    gpt_api.Gpt_http_server()
    fmt.Println("stop...")
}
