package api_from_ai

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"log"
	"strings"
)

func Gpt_api_test() string {
	str := "{\"id\":\"chatcmpl-6yDInI6keCizywdRetjamsS9Qlq3S\",\"object\":\"chat.completion\",\"created\":1679808833,\"model\":\"gpt-3.5-turbo-0301\",\"usage\":{\"prompt_tokens\":15,\"completion_tokens\":112,\"total_tokens\":127},\"choices\":[{\"message\":{\"role\":\"assistant\",\"content\":\"作为一名AI，我没有使用过Java，但从程序员的角度来看，Java有很多好用的特点，如跨平台、面向对象、安全等。Java的语法也易于学习和使用，适合初学者入门。此 外，Java在Web开发、大数据、移动开发等方面也有广泛应用，是一门非常流行和实用的编程语言。\"},\"finish_reason\":\"stop\",\"index\":0}]}"

	jsonReader, _ := simplejson.NewFromReader(strings.NewReader(str))
	content, _ := jsonReader.Get("choices").GetIndex(0).Get("message").Get("content").String()
	log.Println("simplejson:", content)

	/*
	   content := gojson.Json(str).Get("choices").Getindex(1).Get("message").Get("content").Tostring()
	   log.Println(string(content))
	*/
	messagesInfo := []MessagesInfo{
		{
			Role:    "user",
			Content: "java mean?",
		},
	}
	contentInfo := ContentInfo{
		Model:    "gpt-3.5-turbo",
		Messages: messagesInfo,
	}
	/**/
	data, _ := json.Marshal(contentInfo)
	// fmt.Printf("%s", data)
	log.Println(string(data))

	//json.NewEncoder(os.Stdout).Encode(contentInfo)

	//fmt.Println("test")
	return string(data)
}

// POST https://api.openai.com/v1/chat/completions
// Authorization: Bearer sk-xxx
// Content-Type: application/json
// {"model":"gpt-3.5-turbo","messages":[{"role":"user","content":"what is gpt?"}]}
/*
   {
    "model": "gpt-3.5-turbo",
    "messages": [
      {
        "role": "user",
        "content": "what is gpt?"
      }
    ]
   }
*/
