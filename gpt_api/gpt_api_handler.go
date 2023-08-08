package gpt_api

import (
    "encoding/json"
    "github.com/bitly/go-simplejson"
    "chatgpt_api/domain"
    "chatgpt_api/utils"
    "io"
    "log"
    "net/http"
    "strings"
    "time"
)

type MessagesInfo struct {
    Role    string `json:"role""`
    Content string `json:"content"`
}
type ContentInfo struct {
    Model    string         `json:"model"`
    Messages []MessagesInfo `json:"messages"`
}

func GptApi(content string) (string, error) {
    return GptApi2(content, nil)
}

func GptApi2(content string, keywordsArr []domain.Keywords) (string, error) {
    url := "https://api.openai.com/v1/chat/completions"


    log.Printf("  --> ask openai.com ...【%s】 (%s)\n", content, url)

    messagesInfo := []MessagesInfo{}

    for _, keywords := range keywordsArr {
        messagesInfo = append(messagesInfo, MessagesInfo{Role: "user", Content: keywords.Keyword})
        messagesInfo = append(messagesInfo, MessagesInfo{Role: "assistant", Content: keywords.Answer})
    }
    log.Printf("附带历史消息 %d 条", len(keywordsArr))

    messagesInfo = append(messagesInfo, MessagesInfo{ Role: "user", Content: content})

    // 模型版本
    contentInfo := ContentInfo{
        Model:    "gpt-3.5-turbo",
        Messages: messagesInfo,
    }

    //contentInit := "{\"model\":\"gpt-3.5-turbo\",\"messages\":[{\"role\":\"user\",\"content\":\"\"}]}"
    //jsonDecoder := json.NewEncoder(nil)
    //jsonDecoder.Encode(contentInfo)

    httpClient := http.Client{Timeout: 120 * time.Second}

    postContent, _ := json.Marshal(contentInfo)

    req, err := http.NewRequest("POST", url, strings.NewReader(string(postContent)))
    req.Header.Set("Content-Type", "application/json")

    // zxj
    // sk-xxx

    if utils.ApiKey != "" {
        req.Header.Set("Authorization", "Bearer " + utils.ApiKey)
    } else {
        //
        log.Fatal("未配置ApiKey")
    }

    if err != nil {
        log.Println("httpClient NewRequest exception")
        return "", err
    }

    resp, err := httpClient.Do(req)

    if err != nil {
        log.Println("httpClient request exception")
        return "", err
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Println("read exception")
        return "", err
    }

    // {"id":"chatcmpl-6yDInI6keCizywdRetjamsS9Qlq3S","object":"chat.completion","created":1679808833,"model":"gpt-3.5-turbo-0301",
    // "usage":{"prompt_tokens":15,"completion_tokens":112,"total_tokens":127},
    // "choices":[{"message":{"role":"assistant","content":"作为一名AI，我没有使用过Java，但从程序员的角度来看，Java有很多好用的特点，如跨平台、面向对象、安全等。Java的语法也易于学习和使用，适合初学者入门。此 外，Java在Web开发、大数据、移动开发等方面也有广泛应用，是一门非常流行和实用的编程语言。"},"finish_reason":"stop","index":0}]}

    respStr := string(respBody)

    jsonReader, _ := simplejson.NewFromReader(strings.NewReader(respStr))
    respContentStr, _ := jsonReader.Get("choices").GetIndex(0).Get("message").Get("content").String()

    log.Printf("  <-- %s 调用openapi.com完成 %s", content, utils.Substring(respContentStr, 20))
    return respContentStr, nil
}

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
