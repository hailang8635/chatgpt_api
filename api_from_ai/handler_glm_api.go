package api_from_ai

import (
	"bytes"
	"chatgpt_api/config"
	"chatgpt_api/domain"
	"chatgpt_api/utils"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const (
//API_KEY  = "your_api_key_here"
//API_URL  = "https://open.bigmodel.cn/api/paas/v3/model-api/chatglm_turbo/invoke"
//API_KEY = ""
//API_URL = "https://open.bigmodel.cn/api/paas/v4/chat/completions/"
)

/**
 * 返回答案，字符串的格式
 */
func GLMApi(content string) (string, error) {
	return GLMApiWithHistory(content, nil)
}

func GLMApiWithHistory(content string, keywordsArr []domain.KeywordAndAnswerItem) (string, error) {
	log.Printf("  --> ask open.bigmodel.cn ...【%s】 (%s)\n", content, config.GLM_Url)

	messagesInfo := []Message{}

	for _, keywords := range keywordsArr {
		messagesInfo = append(messagesInfo, Message{Role: "user", Content: keywords.Keyword})
		// messagesInfo = append(messagesInfo, Message{Role: "assistant", Content: keywords.Answer})
	}
	log.Printf("附带历史消息 %d 条", len(keywordsArr))

	messagesInfo = append(messagesInfo, Message{Role: "user", Content: content})

	// 构建请求体
	// 或 glm-3-turbo/glm-4
	// Temperature: 0.7,
	// 你好，请用一句话介绍你自己
	// 1+2+...+100等于几
	requestBody := ChatGLMRequest{
		// Model:    "glm-4-plus",
		Model:    config.GLM_Model,
		Messages: messagesInfo,
		//Messages: []Message{{Role: "user", Content: content}},
	}

	jsonBody, _ := json.Marshal(requestBody)

	//fmt.Println("jsonBody request", requestBody)

	// 创建 HTTP 请求
	req, _ := http.NewRequest("POST", config.GLM_Url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", config.GLM_Apikey)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("req json: ", config.GLM_Model, config.GLM_Url, messagesInfo)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// 处理响应
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.Println("StatusCode错误 resp.Status:", resp.Status)
	}

	var result ChatGLMResponse2
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	// fmt.Println("ChatGLMResponse2: ", result)
	if result.Error.Code == "" {
		log.Println("回复:", utils.Substring(result.Choices[0].Message.Content, 100)+"\n\n")
	} else {
		log.Println("错误:", result.Error.Message)
	}

	//respStr := string(respBody)

	//jsonReader, _ := simplejson.NewFromReader(strings.NewReader(respStr))
	//respContentStr, _ := jsonReader.Get("choices").GetIndex(0).Get("message").Get("content").String()
	respContentStr := result.Choices[0].Message.Content

	log.Printf("  <-- %s 调用 open.bigmodel.cn 完成 %s", content, utils.Substring(respContentStr, 20))
	return respContentStr, nil
}

type ChatGLMResponse2 struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Choices []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type ChatGLMRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

/**
{
	"error": {
		"code": "1001",
		"message": "Header中未收到Authorization参数，无法进行身份验证。"
	}
}

{
	"choices": [
		{
			"finish_reason": "stop",
			"index": 0,
			"message": {
				"content": "我是基于人工智能技术的智能助手，致力于为用户提供高效、准确的信息查询和问题解答服务。",
				"role": "assistant"
			}
		}
	],
	"created": 1740475399,
	"id": "20250225172317de4c1f586e1b499e",
	"model": "glm-4-plus",
	"request_id": "20250225172317de4c1f586e1b499e",
	"usage": {
		"completion_tokens": 22,
		"prompt_tokens": 11,
		"total_tokens": 33
	}
}
*/
type ChatGLMResponse struct {
	//Code int    `json:"code"`
	//Msg  string `json:"msg"`
	//Data struct {
	Choices []struct {
		Content string `json:"content"`
	} `json:"choices"`
	//} `json:"data"`
}
