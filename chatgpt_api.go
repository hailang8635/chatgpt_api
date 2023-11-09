package main

import (
    "fmt"
    "chatgpt_api/gpt_api"
)

func main() {

    fmt.Println("start...")
    gpt_api.Gpt_http_server()
    fmt.Println("stop...")
}
